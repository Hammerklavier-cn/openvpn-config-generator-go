package cerinit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

// initialise PKI.
// This is a replacement for `./easyrsa init-pki`.
func initPKI(dir string, verbose bool) error {

	var EASYRSA_PKI = path.Join(dir, "pki")
	// var vars_file = "vars"
	// var vars_file_example = "vars.example"
	var ssl_cnf_file = "openssl-easyrsa.cnf"
	var x509_types_dir = "x509-types"

	// Create PKI directory
	if fileInfo, _ := os.Stat(EASYRSA_PKI); fileInfo != nil {
		fmt.Printf("WARNING: PKI directory already exists! Existing PKI directory will be removed and recreated.\n")
		// TODO: add confirmation before removing the directory here.
		if err := os.RemoveAll(EASYRSA_PKI); err != nil {
			return err
		}
	}
	if err := os.Mkdir(EASYRSA_PKI, 0755); err != nil {
		return err
	}

	// Create subdirectories of `pki`
	for _, area := range []string{"private", "req", "inline"} {
		if err := os.Mkdir(path.Join(EASYRSA_PKI, area), 0755); err != nil {
			return err
		}
	}

	// native implementation of `./easyrsa`'s `install_data_to_pki init-pki` function
	/*
		Note from `./easyrsa`'s `install_data_to_pki:
		# Explicitly find and optionally copy data-files to the PKI.
		# During 'init-pki' this is the new default.
		# During all other functions these requirements are tested for
		# and files will be copied to the PKI, if they do not already
		# exist there.
		#
		# One reason for this is to make packaging work.
	*/
	var areas = []string{
		EASYRSA_PKI,
		".",
		"/usr/local/share/easy-rsa",
		"/usr/share/easy-rsa",
		"/etc/easy-rsa",
	}

	var EasyrsaExtDir string

	for _, area := range areas {

		// Find x509-types and keep the first one found
		if fileStat, _ := os.Stat(path.Join(area, x509_types_dir)); fileStat != nil {
			if fileStat.IsDir() && EasyrsaExtDir == "" {
				EasyrsaExtDir = path.Join(area, "x509-types")
			}
		}

		// find 'openssl-easyrsa.cnf'
		if fileStat, _ := os.Stat(path.Join(area, ssl_cnf_file)); fileStat != nil {
			if _, err := os.Stat(path.Join(EASYRSA_PKI, ssl_cnf_file)); os.IsNotExist(err) {
				src, err := os.Open(path.Join(area, ssl_cnf_file))
				if err != nil {
					return err
				}
				defer src.Close()

				dst, err := os.Create(path.Join(EASYRSA_PKI, ssl_cnf_file))
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err := io.Copy(dst, src); err != nil {
					return err
				}
			} else {
				continue
			}
		} else {
			continue
		}
	}

	if EasyrsaExtDir == "" {
		return errors.New("EasyrsaExtDir not found")
	} else {
		fmt.Println("EasyrsaExtDir found in", EasyrsaExtDir)
	}

	// return err if 'openssl-easyrsa.cnf' not found
	if _, err := os.Stat(path.Join(EASYRSA_PKI, ssl_cnf_file)); os.IsNotExist(err) {
		return errors.New("openssl-easyrsa.cnf not found")
	}

	// The following commented lines are deprecated.
	// The original bash code of this part doesn't make any
	// sense under this circumstance, because `vars.example`
	// does not exists under the PKI directory.
	//
	// // create `vars` if not found
	// if _, err := os.Stat(path.Join(EASYRSA_PKI, vars_file)); os.IsNotExist(err) {
	// 	if file_stat, _ := os.Stat(path.Join(EASYRSA_PKI, vars_file_example)); file_stat != nil && !file_stat.IsDir() {
	// 		source_file, err := os.Open(path.Join(EASYRSA_PKI, vars_file_example))
	// 		if err != nil {
	// 			return err
	// 		}
	// 		defer source_file.Close()

	// 		dest_file, err := os.Create(path.Join(EASYRSA_PKI, vars_file))
	// 		if err != nil {
	// 			return err
	// 		}
	// 		defer dest_file.Close()

	// 		if _, err := io.Copy(dest_file, source_file); err != nil {
	// 			return err
	// 		}

	// 		fmt.Println("vars is created under", path.Join(EASYRSA_PKI, vars_file))
	// 	} else {
	// 		fmt.Printf(
	// 			"vars.example not found at %s. Please create vars file manually.\n",
	// 			path.Join(EASYRSA_PKI, vars_file_example))
	// 	}
	// } else {
	// 	if verbose {
	// 		fmt.Printf(
	// 			"vars is found under %s. Skip creating new one.\n",
	// 			path.Join(EASYRSA_PKI, vars_file))
	// 	}
	// }

	return nil
}

// This is a replacement for `./easyrsa build-ca`.
//
// It will generate ca.crt and private/ca.key under pki directory.
func buildCA(dir string, verbose bool) error {

	var EASYRSA_PKI = path.Join(dir, "pki")
	// var EASYRSA_REQ_CN = "Easy-RSA CA"

	// var cipher = "-aes256"
	// var nopass = true
	// var out_key = path.Join(EASYRSA_PKI, "private", "ca.key")
	// var out_file = path.Join(EASYRSA_PKI, "ca.crt")
	// var date_stamp = 1
	// var x509 = 1

	// if nopass {
	// 	cipher = ""
	// }

	// The following is the go implementation of `verify_ca_init test`
	// function of `easyrsa`
	{
		// Check if any of the following files exists
		file_names := []string{
			"ca.crt", path.Join("private", "ca.key"),
			"index.txt", "index.txt.attr", "serial"}
		for _, file_name := range file_names {
			if _, err := os.Stat(path.Join(EASYRSA_PKI, file_name)); err == nil {
				return fmt.Errorf(
					"Found existing CA file %s here, which is unexpected.\n",
					path.Join(EASYRSA_PKI, file_name))
			}
		}
	}
	// `verify_ca_init test` implementation ends.

	// create necessary dirs and files
	{
		// create necessary dirs
		dirs := []string{
			"issued", "certs_by_serial",
			path.Join("revoked", "certs_by_serial"),
			path.Join("revoked", "private_by_serial"),
			path.Join("revoked", "reqs_by_serial"),
		}
		for _, dir := range dirs {
			if err := os.MkdirAll(path.Join(EASYRSA_PKI, dir), 0755); err != nil {
				return err
			}
		}

		// create necessary files
		file, err := os.OpenFile(path.Join(EASYRSA_PKI, "index.txt"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(""); err != nil {
			return err
		}

		file, err = os.OpenFile(path.Join(EASYRSA_PKI, "index.txt.attr"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(fmt.Sprintf("%s\n", "unique_subject = no")); err != nil {
			return err
		}

		file, err = os.OpenFile(path.Join(EASYRSA_PKI, "serial"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(fmt.Sprintf("%s\n", "01")); err != nil {
			return err
		}
	}

	// Assign cert and key temp files
	//
	// The following is the implementation of
	// "easyrsa easyrsa_mktemp()` function

	return nil
}
