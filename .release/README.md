# Release signing material

- `gpg-public.asc` — public key to paste into registry.terraform.io when registering this provider
- `gpg-fingerprint.txt` — key fingerprint (`22D75718E471A28B5EAB9E3139D1567099DBF1DA`)
- `gpg-passphrase.local` — **gitignored**; local copy of the passphrase also stored as the `PASSPHRASE` GitHub Actions secret

Private key lives only in the `GPG_PRIVATE_KEY` GitHub Actions secret (and the local GnuPG homedir `~/.gnupg-fireweave-tf`).
