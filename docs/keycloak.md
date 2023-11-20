---
layout: default
title: Keycloak lib
nav_order: 3
description: "HWM Go Utils Docs"
permalink: /keycloak
last_modified_date: 2023-11-16T12:00:00+0000
---

# Keycloak utils

This section of the HWM Go Utils library provides a set of utility functions for interacting with the Keycloak Identity and Access Management system. These utilities simplify the process of integrating Keycloak into your Go applications.

## Features

### Authentication
This module provides functions for authenticating users with Keycloak. It includes functions for generating and validating tokens, refreshing tokens, and handling token expiration.

### User Management
This module provides functions for managing users in Keycloak. It includes functions for creating, updating, and deleting users, as well as for managing user roles and permissions.

### Realm Management
This module provides functions for managing realms in Keycloak. It includes functions for creating, updating, and deleting realms, as well as for managing realm settings and clients.

## Usage

To use the Keycloak utilities, first import the package:

```go
import "github.com/HiWay-Media/hwm-go-utils/keycloak"
```
