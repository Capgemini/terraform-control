# Terraform-Control

[![Go Report Card](https://goreportcard.com/badge/github.com/Capgemini/terraform-control)](https://goreportcard.com/report/github.com/Capgemini/terraform-control)
[![wercker status](https://app.wercker.com/status/15890eedfe803a8ff9d9044261c0eea7/m "wercker status")](https://app.wercker.com/project/bykey/15890eedfe803a8ff9d9044261c0eea7)
[![Coverage Status](https://coveralls.io/repos/github/Capgemini/terraform-control/badge.svg?branch=HEAD)](https://coveralls.io/github/Capgemini/terraform-control?branch=HEAD)
[![Code Climate](https://codeclimate.com/github/Capgemini/terraform-control/badges/gpa.svg)](https://codeclimate.com/github/Capgemini/terraform-control)


Terraform-Control is a solution for managing and deploying your infrastructure with terraform in a collaborative way driven by continuous integration while keeping track of the state and history of your infrastructure.

## Overview

We have reused loads of the [Otto](https://github.com/hashicorp/otto/) code for running terraform over different environments simulating an [Atlas terraform](https://atlas.hashicorp.com/terraform) style solution to demonstrate how to use terraform in a collaborative way driven by continuous integration while keeping track of the state of your environment in a centralised way.
**This is just a PoC and it's obviously missing a lot features to be used in a real environment at the minute.**

![terraform-control-diagram](docs/terraform-control-diagram.png)

## Demo

[![Terraform-control PoC](https://img.youtube.com/vi/5eClxFWK_Ec/0.jpg)](https://www.youtube.com/watch?v=5eClxFWK_Ec)


## Web UI

![web-ui](docs/terraform-control-ui.gif)

## Blog

https://capgemini.github.io/devops/Controlling-the-state-of-your-infrastructure/
