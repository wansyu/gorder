import base64
import hashlib

import scrypt

def main():
    # DO NOT use this salt value; generate your own random salt. 8 bytes is
    # a good length.
    salt = '@#$%^&*()'
    dk = scrypt.hash("somepassword", salt, N=2**15, r=8, p=1, buflen=32)
    print(base64.b64encode(dk).decode())

if __name__ == "__main__":
    main()
