---
page_title: "Fastly: service_waf_configuration"
---

-> **Note:** This resource is only available from 0.20.0 of the Fastly terraform provider.

# fastly_service_waf_configuration

Defines a set of Web Application Firewall configuration options that can be used to populate a service WAF. This resource will configure rules, thresholds and other settings for a WAF.


~> **Warning:** Terraform will take precedence over any changes you make in the UI or API. Such changes are likely to be reversed if you run Terraform again.


## Example Usage

Basic usage:

```hcl
resource "fastly_service_v1" "demo" {
  name = "demofastly"

  domain {
    name    = "example.com"
    comment = "demo"
  }

  backend {
    address = "127.0.0.1"
    name    = "origin1"
    port    = 80
  }

  condition {
    name      = "WAF_Prefetch"
    type      = "PREFETCH"
    statement = "req.backend.is_origin"
  }

  # This condition will always be false
  # adding it to the response object created below
  # prevents Fastly from returning a 403 on all of your traffic.
  condition {
    name      = "WAF_always_false"
    statement = "false"
    type      = "REQUEST"
  }

  response_object {
    name              = "WAF_Response"
    status            = "403"
    response          = "Forbidden"
    content_type      = "text/html"
    content           = "<html><body>Forbidden</body></html>"
    request_condition = "WAF_always_false"
  }

  waf {
    prefetch_condition = "WAF_Prefetch"
    response_object    = "WAF_Response"
  }

  force_destroy = true
}

resource "fastly_service_waf_configuration" "waf" {
  waf_id                          = fastly_service_v1.demo.waf[0].waf_id
  http_violation_score_threshold  = 100
}
```

Usage with rules:

```hcl
resource "fastly_service_v1" "demo" {
  name = "demofastly"

  domain {
    name    = "example.com"
    comment = "demo"
  }

  backend {
    address = "127.0.0.1"
    name    = "origin1"
    port    = 80
  }

  condition {
    name      = "WAF_Prefetch"
    type      = "PREFETCH"
    statement = "req.backend.is_origin"
  }

  # This condition will always be false
  # adding it to the response object created below
  # prevents Fastly from returning a 403 on all of your traffic.
  condition {
    name      = "WAF_always_false"
    statement = "false"
    type      = "REQUEST"
  }

  response_object {
    name              = "WAF_Response"
    status            = "403"
    response          = "Forbidden"
    content_type      = "text/html"
    content           = "<html><body>Forbidden</body></html>"
    request_condition = "WAF_always_false"
  }

  waf {
    prefetch_condition = "WAF_Prefetch"
    response_object    = "WAF_Response"
  }

  force_destroy = true
}

resource "fastly_service_waf_configuration" "waf" {
  waf_id                          = fastly_service_v1.demo.waf[0].waf_id
  http_violation_score_threshold  = 100

  rule {
    modsec_rule_id = 1010090
    revision       = 1
    status         = "log"
  }
}
```

Usage with rule exclusions:

~> **Warning:** Rule exclusions are part of a **beta release**, which may be subject to breaking changes and improvements over time. For more information, see our [product and feature lifecycle](https://docs.fastly.com/products/fastly-product-lifecycle#beta) descriptions.

```hcl
resource "fastly_service_v1" "demo" {
  name = "demofastly"

  domain {
    name    = "example.com"
    comment = "demo"
  }

  backend {
    address = "127.0.0.1"
    name    = "origin1"
    port    = 80
  }

  condition {
    name      = "WAF_Prefetch"
    type      = "PREFETCH"
    statement = "req.backend.is_origin"
  }

  # This condition will always be false
  # adding it to the response object created below
  # prevents Fastly from returning a 403 on all of your traffic.
  condition {
    name      = "WAF_always_false"
    statement = "false"
    type      = "REQUEST"
  }

  response_object {
    name              = "WAF_Response"
    status            = "403"
    response          = "Forbidden"
    content_type      = "text/html"
    content           = "<html><body>Forbidden</body></html>"
    request_condition = "WAF_always_false"
  }

  waf {
    prefetch_condition = "WAF_Prefetch"
    response_object    = "WAF_Response"
  }

  force_destroy = true
}

resource "fastly_service_waf_configuration" "waf" {
  waf_id                          = fastly_service_v1.demo.waf[0].waf_id
  http_violation_score_threshold  = 100

  rule {
    modsec_rule_id = 2029718
    revision       = 1
    status         = "log"
  }

  rule_exclusion {
    name = "index page"
    exclusion_type = "rule"
    condition = "req.url.basename == \"index.html\""
    modsec_rule_ids = [2029718]
  }
}
```

Usage with rules from data source:

```hcl
variable "type_status" {
  type = map(string)
  default = {
    score     = "score"
    threshold = "log"
    strict    = "log"
  }
}

resource "fastly_service_v1" "demo" {
  name = "demofastly"

  domain {
    name    = "example.com"
    comment = "demo"
  }

  backend {
    address = "127.0.0.1"
    name    = "origin1"
    port    = 80
  }

  condition {
    name      = "WAF_Prefetch"
    type      = "PREFETCH"
    statement = "req.backend.is_origin"
  }

  # This condition will always be false
  # adding it to the response object created below
  # prevents Fastly from returning a 403 on all of your traffic.
  condition {
    name      = "WAF_always_false"
    statement = "false"
    type      = "REQUEST"
  }

  response_object {
    name              = "WAF_Response"
    status            = "403"
    response          = "Forbidden"
    content_type      = "text/html"
    content           = "<html><body>Forbidden</body></html>"
    request_condition = "WAF_always_false"
  }
  waf {
    prefetch_condition = "WAF_Prefetch"
    response_object    = "WAF_Response"
  }

  force_destroy = true
}

data "fastly_waf_rules" "owasp" {
  publishers = ["owasp"]
}

resource "fastly_service_waf_configuration" "waf" {
  waf_id                          = fastly_service_v1.demo.waf[0].waf_id
  http_violation_score_threshold  = 100

  dynamic "rule" {
    for_each = data.fastly_waf_rules.owasp.rules
    content {
      modsec_rule_id = rule.value.modsec_rule_id
      revision       = rule.value.latest_revision_number
      status         = lookup(var.type_status, rule.value.type, "log")
    }
  }
}
```

Usage with support for individual rule configuration (this is the suggested pattern):

```hcl
# this variable is used for rule configuration in bulk
variable "type_status" {
  type = map(string)
  default = {
    score     = "score"
    threshold = "log"
    strict    = "log"
  }
}
# this variable is used for individual rule configuration
variable "individual_rules" {
  type = map(string)
  default = {
    1010020 = "block"
  }
}

resource "fastly_service_v1" "demo" {
  name = "demofastly"

  domain {
    name    = "example.com"
    comment = "demo"
  }

  backend {
    address = "127.0.0.1"
    name    = "origin1"
    port    = 80
  }

  condition {
    name      = "WAF_Prefetch"
    type      = "PREFETCH"
    statement = "req.backend.is_origin"
  }

  # This condition will always be false
  # adding it to the response object created below
  # prevents Fastly from returning a 403 on all of your traffic.
  condition {
    name      = "WAF_always_false"
    statement = "false"
    type      = "REQUEST"
  }

  response_object {
    name              = "WAF_Response"
    status            = "403"
    response          = "Forbidden"
    content_type      = "text/html"
    content           = "<html><body>Forbidden</body></html>"
    request_condition = "WAF_always_false"
  }

  waf {
    prefetch_condition = "WAF_Prefetch"
    response_object    = "WAF_Response"
  }

  force_destroy = true
}

data "fastly_waf_rules" "owasp" {
  publishers = ["owasp"]
}

resource "fastly_service_waf_configuration" "waf" {
  waf_id                          = fastly_service_v1.demo.waf[0].waf_id
  http_violation_score_threshold  = 202

  dynamic "rule" {
    for_each = data.fastly_waf_rules.owasp.rules
    content {
      modsec_rule_id = rule.value.modsec_rule_id
      revision       = rule.value.latest_revision_number
      # Nested lookups in order to apply a combination of in bulk and individual rule configuration.
      status         = lookup(var.individual_rules, rule.value.modsec_rule_id, lookup(var.type_status, rule.value.type, "log"))
    }
  }
}
```

Usage with support for specific rule revision configuration:

```hcl
# this variable is used for rule configuration in bulk
variable "type_status" {
  type = map(string)
  default = {
    score     = "score"
    threshold = "log"
    strict    = "log"
  }
}

# This variable is used for individual rule revision configuration.
variable "specific_rule_revisions" {
  type = map(string)
  default = {
    #  If the revision requested is not found, the server will return a 404 response code.
    1010020 = 1
  }
}

resource "fastly_service_v1" "demo" {
  name = "demofastly"

  domain {
    name    = "example.com"
    comment = "demo"
  }

  backend {
    address = "127.0.0.1"
    name    = "origin1"
    port    = 80
  }

  condition {
    name      = "WAF_Prefetch"
    type      = "PREFETCH"
    statement = "req.backend.is_origin"
  }

  # This condition will always be false
  # adding it to the response object created below
  # prevents Fastly from returning a 403 on all of your traffic.
  condition {
    name      = "WAF_always_false"
    statement = "false"
    type      = "REQUEST"
  }

  response_object {
    name              = "WAF_Response"
    status            = "403"
    response          = "Forbidden"
    content_type      = "text/html"
    content           = "<html><body>Forbidden</body></html>"
    request_condition = "WAF_always_false"
  }

  waf {
    prefetch_condition = "WAF_Prefetch"
    response_object    = "WAF_Response"
  }

  force_destroy = true
}

data "fastly_waf_rules" "owasp" {
  publishers = ["owasp"]
}

resource "fastly_service_waf_configuration" "waf" {
  waf_id                          = fastly_service_v1.demo.waf[0].waf_id
  http_violation_score_threshold  = 202

  dynamic "rule" {
    for_each = data.fastly_waf_rules.owasp.rules
    content {
      modsec_rule_id = rule.value.modsec_rule_id
      revision       = lookup(var.specific_rule_revisions, rule.value.modsec_rule_id, rule.value.latest_revision_number)
      status         = lookup(var.type_status, rule.value.type, "log")
    }
  }
}
```

Usage omitting rule revision field. The first time Terraform is applied, the latest rule revisions are associated with the WAF. Any subsequent apply would not alter the rule revisions.

```hcl
# This variable is used for rule configuration in bulk.
variable "type_status" {
  type = map(string)
  default = {
    score     = "score"
    threshold = "log"
    strict    = "log"
  }
}
# This variable is used for individual rule configuration.
variable "individual_rules" {
  type = map(string)
  default = {
    1010020 = "block"
  }
}

resource "fastly_service_v1" "demo" {
  name = "demofastly"

  domain {
    name    = "example.com"
    comment = "demo"
  }

  backend {
    address = "127.0.0.1"
    name    = "origin1"
    port    = 80
  }

  condition {
    name      = "WAF_Prefetch"
    type      = "PREFETCH"
    statement = "req.backend.is_origin"
  }

  # This condition will always be false
  # adding it to the response object created below
  # prevents Fastly from returning a 403 on all of your traffic.
  condition {
    name      = "WAF_always_false"
    statement = "false"
    type      = "REQUEST"
  }

  response_object {
    name              = "WAF_Response"
    status            = "403"
    response          = "Forbidden"
    content_type      = "text/html"
    content           = "<html><body>Forbidden</body></html>"
    request_condition = "WAF_always_false"
  }

  waf {
    prefetch_condition = "WAF_Prefetch"
    response_object    = "WAF_Response"
  }

  force_destroy = true
}

data "fastly_waf_rules" "owasp" {
  publishers = ["owasp"]
}

resource "fastly_service_waf_configuration" "waf" {
  waf_id                          = fastly_service_v1.demo.waf[0].waf_id
  http_violation_score_threshold  = 202

  dynamic "rule" {
    for_each = data.fastly_waf_rules.owasp.rules
    content {
      modsec_rule_id = rule.value.modsec_rule_id
      # Rule revision field ommitted.
      status         = lookup(var.individual_rules, rule.value.modsec_rule_id, lookup(var.type_status, rule.value.type, "log"))
    }
  }
}
# This output contains the WAF's active rules set. This can be useful for identifying which revision is active for each WAF rule.
output "rules" {
  value = fastly_service_waf_configuration.waf.rule
}
```

## Adding a WAF to an existing service

~> **Warning:** A two-phase change is required when adding a WAF to an existing service

When adding a `waf` to an existing `fastly_service_v1` and at the same time adding a `fastly_service_waf_configuration`
resource with `waf_id = fastly_service_v1.demo.waf[0].waf_id` might result with the in the following error:

> fastly_service_v1.demo.waf is empty list of object

For this scenario, it's recommended to split the changes into two distinct steps:

1. Add the `waf` block to the `fastly_service_v1` and apply the changes
2. Add the `fastly_service_waf_configuration` to the HCL and apply the changes

## Argument Reference

The following arguments are supported:

* `waf_id` - (Required) The ID of the Web Application Firewall that the configuration belongs to.
* `allowed_http_versions` - (Optional) Allowed HTTP versions.
* `allowed_methods` - (Optional) A space-separated list of HTTP method names.
* `allowed_request_content_type` - (Optional) Allowed request content types.
* `allowed_request_content_type_charset` - (Optional) Allowed request content type charset.
* `arg_length` - (Optional) The maximum number of arguments allowed.
* `arg_name_length` - (Optional) The maximum allowed argument name length.
* `combined_file_sizes` - (Optional) The maximum allowed size of all files.
* `critical_anomaly_score` - (Optional) Score value to add for critical anomalies.
* `crs_validate_utf8_encoding` - (Optional) CRS validate UTF8 encoding.
* `error_anomaly_score` - (Optional) Score value to add for error anomalies.
* `high_risk_country_codes` - (Optional) A space-separated list of country codes in ISO 3166-1 (two-letter) format.
* `http_violation_score_threshold` - (Optional) HTTP violation threshold.
* `inbound_anomaly_score_threshold` - (Optional) Inbound anomaly threshold.
* `lfi_score_threshold` - (Optional) Local file inclusion attack threshold.
* `max_file_size` - (Optional) The maximum allowed file size, in bytes.
* `max_num_args` - (Optional) The maximum number of arguments allowed.
* `notice_anomaly_score` - (Optional) Score value to add for notice anomalies.
* `paranoia_level` - (Optional) The configured paranoia level.
* `php_injection_score_threshold` - (Optional) PHP injection threshold.
* `rce_score_threshold` - (Optional) Remote code execution threshold.
* `restricted_extensions` - (Optional) A space-separated list of allowed file extensions.
* `restricted_headers` - (Optional) A space-separated list of allowed header names.
* `rfi_score_threshold` - (Optional) Remote file inclusion attack threshold.
* `session_fixation_score_threshold` - (Optional) Session fixation attack threshold.
* `sql_injection_score_threshold` - (Optional) SQL injection attack threshold.
* `total_arg_length` - (Optional) The maximum size of argument names and values.
* `warning_anomaly_score` - (Optional) Score value to add for warning anomalies.
* `xss_score_threshold` - (Optional) XSS attack threshold.
* `rule` - (Optional) The Web Application Firewall's active rules. [Defined below](#rule-block)
* `rule_exclusion` - (Optional) The Web Application Firewall's rule exclusions. [Defined below](#rule_exclusion-block)

### rule block

The `rule` block supports:

* `status` - (Required) The Web Application Firewall rule's status. Allowed values are (`log`, `block` and `score`).
* `modsec_rule_id` - (Required) The Web Application Firewall rule's modsecurity ID.
* `revision` - (Optional) The Web Application Firewall rule's revision. The latest revision will be used if this is not provided.

### rule_exclusion block

The `rule_exclusion` block supports:

~> **Warning:** Rule exclusions are part of a **beta release**, which may be subject to breaking changes and improvements over time. For more information, see our [product and feature lifecycle](https://docs.fastly.com/products/fastly-product-lifecycle#beta) descriptions.

* `name` - (Required) The name of rule exclusion.
* `exclusion_type` - (Required) The type of rule exclusion. Values are `rule` to exclude the specified rule(s), or `waf` to disable the Web Application Firewall.
* `condition` - (Required) A conditional expression in VCL used to determine if the condition is met.
* `modsec_rule_ids` - (Required) Set of modsecurity IDs to be excluded. No rules should be provided when `exclusion_type` is `waf`. The rules need to be configured on the Web Application Firewall to be excluded.
* `number` - The numeric ID assigned to the WAF Rule Exclusion.

## Import

This is an example of the import command being applied to the resource named `fastly_service_waf_configuration.waf`
The resource ID should be the WAF ID.

```
$ terraform import fastly_service_waf_configuration.waf xxxxxxxxxxxxxxxxxxxx
```

If Terraform is already managing a remote WAF configurations against a resource being imported then the user will be asked to remove it from the existing Terraform state.
The following is an example of the Terraform state command to remove the resource named `fastly_service_waf_configuration.waf` from the Terraform state file.

```
$ terraform state rm fastly_service_waf_configuration.waf
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **waf_id** (String) The service the WAF belongs to.

### Optional

- **allowed_http_versions** (String) Allowed HTTP versions (default HTTP/1.0 HTTP/1.1 HTTP/2).
- **allowed_methods** (String) A space-separated list of HTTP method names (default GET HEAD POST OPTIONS PUT PATCH DELETE).
- **allowed_request_content_type** (String) Allowed request content types (default application/x-www-form-urlencoded|multipart/form-data|text/xml|application/xml|application/x-amf|application/json|text/plain).
- **allowed_request_content_type_charset** (String) Allowed request content type charset (default utf-8|iso-8859-1|iso-8859-15|windows-1252).
- **arg_length** (Number) The maximum number of arguments allowed (default 400).
- **arg_name_length** (Number) The maximum allowed argument name length (default 100).
- **combined_file_sizes** (Number) The maximum allowed size of all files (in bytes, default 10000000).
- **critical_anomaly_score** (Number) Score value to add for critical anomalies (default 6).
- **crs_validate_utf8_encoding** (Boolean) CRS validate UTF8 encoding.
- **error_anomaly_score** (Number) Score value to add for error anomalies (default 5).
- **high_risk_country_codes** (String) A space-separated list of country codes in ISO 3166-1 (two-letter) format.
- **http_violation_score_threshold** (Number) HTTP violation threshold.
- **id** (String) The ID of this resource.
- **inbound_anomaly_score_threshold** (Number) Inbound anomaly threshold.
- **lfi_score_threshold** (Number) Local file inclusion attack threshold.
- **max_file_size** (Number) The maximum allowed file size, in bytes (default 10000000).
- **max_num_args** (Number) The maximum number of arguments allowed (default 255).
- **notice_anomaly_score** (Number) Score value to add for notice anomalies (default 4).
- **paranoia_level** (Number) The configured paranoia level (default 1).
- **php_injection_score_threshold** (Number) PHP injection threshold.
- **rce_score_threshold** (Number) Remote code execution threshold.
- **restricted_extensions** (String) A space-separated list of allowed file extensions (default .asa/ .asax/ .ascx/ .axd/ .backup/ .bak/ .bat/ .cdx/ .cer/ .cfg/ .cmd/ .com/ .config/ .conf/ .cs/ .csproj/ .csr/ .dat/ .db/ .dbf/ .dll/ .dos/ .htr/ .htw/ .ida/ .idc/ .idq/ .inc/ .ini/ .key/ .licx/ .lnk/ .log/ .mdb/ .old/ .pass/ .pdb/ .pol/ .printer/ .pwd/ .resources/ .resx/ .sql/ .sys/ .vb/ .vbs/ .vbproj/ .vsdisco/ .webinfo/ .xsd/ .xsx).
- **restricted_headers** (String) A space-separated list of allowed header names (default /proxy/ /lock-token/ /content-range/ /translate/ /if/).
- **rfi_score_threshold** (Number) Remote file inclusion attack threshold.
- **rule** (Block Set) (see [below for nested schema](#nestedblock--rule))
- **rule_exclusion** (Block Set) (see [below for nested schema](#nestedblock--rule_exclusion))
- **session_fixation_score_threshold** (Number) Session fixation attack threshold.
- **sql_injection_score_threshold** (Number) SQL injection attack threshold.
- **total_arg_length** (Number) The maximum size of argument names and values (default 6400).
- **warning_anomaly_score** (Number) Score value to add for warning anomalies.
- **xss_score_threshold** (Number) XSS attack threshold.

<a id="nestedblock--rule"></a>
### Nested Schema for `rule`

Required:

- **modsec_rule_id** (Number) The Web Application Firewall rule's modsec ID.
- **status** (String) The Web Application Firewall rule's status. Allowed values are (log, block and score).

Optional:

- **revision** (Number) The Web Application Firewall rule's revision.


<a id="nestedblock--rule_exclusion"></a>
### Nested Schema for `rule_exclusion`

Required:

- **condition** (String) A conditional expression in VCL used to determine if the condition is met.
- **exclusion_type** (String) The type of exclusion.
- **name** (String) Name of the exclusion.

Optional:

- **modsec_rule_ids** (Set of Number) The modsec rule IDs to exclude.

Read-only:

- **number** (Number) A sequential ID assigned to the exclusion.