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
func initPKI(dir string) error {

	// Create PKI directory
	if fileInfo, _ := os.Stat(path.Join(dir, "pki")); fileInfo != nil {
		fmt.Printf("WARNING: PKI directory already exists! Existing PKI directory will be removed and recreated.\n")
		if err := os.RemoveAll(path.Join(dir, "pki")); err != nil {
			return err
		}
	}
	if err := os.Mkdir(path.Join(dir, "pki"), 0755); err != nil {
		return err
	}

	// Create subdirectories of `pki`
	for _, area := range []string{"private", "req", "inline"} {
		if err := os.Mkdir(path.Join(dir, "pki", area), 0755); err != nil {
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
		path.Join(dir, "pki"),
		".",
		"/usr/local/share/easy-rsa",
		"/usr/share/easy-rsa",
		"/etc/easy-rsa",
	}

	var EasyrsaExtDir string

	for _, area := range areas {

		// Find x509-types and keep the first one found
		if fileStat, _ := os.Stat(area); fileStat != nil {
			if fileStat.IsDir() && EasyrsaExtDir == "" {
				EasyrsaExtDir = area
			}
		}

		// find 'openssl-easyrsa.cnf'
		if fileStat, _ := os.Stat(path.Join(area, "openssl-easyrsa.cnf")); fileStat != nil {
			if _, err := os.Stat(path.Join(dir, "pki", "openssl-easyrsa.cnf")); os.IsNotExist(err) {
				src, err := os.Open(path.Join(area, "openssl-easyrsa.cnf"))
				if err != nil {
					return err
				}
				defer src.Close()

				dst, err := os.Create(path.Join(dir, "pki", "openssl-easyrsa.cnf"))
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
	// TODO: return err if EasyrsaExtDir == ""
	if EasyrsaExtDir == "" {
		return errors.New("EasyrsaExtDir not found")
	}

	// TODO: return err if 'openssl-easyrsa.cnf' not found
	if _, err := os.Stat(path.Join(dir, "pki", "openssl-easyrsa.cnf")); os.IsNotExist(err) {
		return errors.New("openssl-easyrsa.cnf not found")
	}

	return nil
}
