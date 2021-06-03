
![](img/lucha128x128.png)

# lucha

![GitHub release (latest by date)](https://img.shields.io/github/v/release/devops-kung-fu/lucha) [![Go Report Card](https://goreportcard.com/badge/github.com/devops-kung-fu/lucha)](https://goreportcard.com/report/github.com/devops-kung-fu/lucha) [![codecov](https://codecov.io/gh/devops-kung-fu/lucha/branch/main/graph/badge.svg?token=P9WBOBQTOB)](https://codecov.io/gh/devops-kung-fu/lucha) [![SBOM](https://img.shields.io/badge/CyloneDX-SBoM-informational)](lucha-sbom.json)

A CLI that scans for sensitive data in source code

## Overview

If you are scanning for secrets with a GitHub Action on a Pull Request then you're way too late. Your secrets are already in your repository and have polluted your history. Rather than commit secrets in the first place, we developed ```lucha``` to root them out - and when combined with [Hookz](https://github.com/devops-kung-fu/hookz), commits can fail before they are pushed to your remote repository. The best part? This works with any ```git``` based remote as everything happens locally on your machine.

Talk about shifting left, right?

## Secret Detection

```lucha``` contains a number of rules that can detect secrets, keys, and tokens that exist in your codebase. The following list are a number of secrets that ```lucha``` can find:

* Github Personal Access Tokens
* AWS Client Identifiers
* AWS Access Keys
* Mailgun API Keys
* Twitter Access Tokens
* Facebook Access Tokens
* Google API Keys
* Stripe API Keys
* Square Access Tokens
* Square OAuth Secrets
* PayPal/Braintree Access Tokens
* AMS Auth Tokens
* MailChimp API Keys
* Slack API Keys
* Slack Access Tokens
* GCP OAuth2.0 Credentials
* GCP API Keys

## Installation

TBD


## Software Bill of Materials

```lucha``` uses [Hookz](https://github.com/devops-kung-fu/hookz) and CycloneDX to generate a Software Bill of Materials in CycloneDX format every time a developer commits code to this repository. More information for CycloneDX is available [here](https://cyclonedx.org)

The current SBoM for ```lucha``` is available [here](lucha-sbom.json).

## Credits

A big thank-you to our friends at [Freepik](https://www.freepik.com) for the ```lucha``` logo.
