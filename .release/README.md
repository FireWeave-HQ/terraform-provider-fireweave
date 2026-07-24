# Release signing (maintainers)

Material used to sign Terraform Registry releases for this provider.

| File | Purpose |
|------|---------|
| `gpg-public.asc` | Public key registered with the Terraform Registry |
| `gpg-fingerprint.txt` | Key fingerprint |
| `gpg-passphrase.local` | **Gitignored.** Local passphrase copy; also stored as the `PASSPHRASE` Actions secret |

The private key is stored only as the `GPG_PRIVATE_KEY` GitHub Actions secret (and optionally in a local GnuPG homedir for operators). Never commit private key material.

**Fingerprint:** `22D75718E471A28B5EAB9E3139D1567099DBF1DA`
