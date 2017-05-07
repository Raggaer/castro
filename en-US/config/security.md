---
name: Security
---

Provides access to the security layer of Castro. Most of these values should never be touched.

- [XSS](#xss)
- [Frame](#frame)
- [ContentType](#contenttype)
- [ReferrerPolicy](#referrerpolicy)
- [CrossDomainPolicy](#crossdomainpolicy)
- [STS](#sts)
- [Security.CSP](#csp)

# XSS

Sets the XSS security header. By default the value is `1; mode=block`.

# Frame

Sets the Frame security header. By default the value is ` DENY`. This will not allow to get Castro as an `iframe`.

# ContentType

Makes sure all the MIME headers are followed and do not change. By default the value is `nosniff`.

# ReferrerPolicy

Sets the referrer policy header. By default the value is `origin`.

# CrossDomainPolicy

Specifies the cross domain policy file to use. By default the value is `none`.

# STS

Specifies the amount of seconds the browser should cache your SSL cert. By default the value is `max-age=10000`.

# CSP

Provides access to Castro Content Security Policy header.

- [CSP.Enabled](#csp.enabled)
- [CSP.Default](#csp.default)
- [CSP.Frame](#csp.frame)

# CSP.Enabled

Specifies if Content Security Policy is enabled. By default this value is `true`.

# CSP.Default

Sets the default values for the policy. This value is an array of strings. By default the value is `none`.

# CSP.Frame

Provides access to the Frame part of the Content Security Policy. This layer contains two fields:

- [Frame.Default](#frame.default)
- [Frame.SRC](#frame.src)

# Frame.Default

Sets the default values for the policy. This value is an array of strings.

# Frame.SRC

Sets the SRC values for the policy. This value is an array of strings.

# CSP.Script

Provides access to the Script part of the Content Security Policy. This layer contains two fields:

- [Script.Default](#script.default)
- [Script.SRC](#script.src)

# Script.Default

Sets the default values for the policy. This value is an array of strings.

# Script.SRC

Sets the SRC values for the policy. This value is an array of strings.

# CSP.Font

Provides access to the Font part of the Content Security Policy. This layer contains two fields:

- [Font.Default](#font.default)
- [Font.SRC](#font.src)

# Font.Default

Sets the default values for the policy. This value is an array of strings.

# Font.SRC

Sets the SRC values for the policy. This value is an array of strings.
 
# CSP.Connect

Provides access to the Connect part of the Content Security Policy. This layer contains two fields:

- [Connect.Default](#connect.default)
- [Connect.SRC](#connect.src)

# Connect.Default

Sets the default values for the policy. This value is an array of strings.

# Connect.SRC

Sets the SRC values for the policy. This value is an array of strings.

# CSP.Style

Provides access to the Style part of the Content Security Policy. This layer contains two fields:

- [Style.Default](#style.default)
- [Style.SRC](#style.src)

# Style.Default

Sets the default values for the policy. This value is an array of strings.

# Style.SRC

Sets the SRC values for the policy. This value is an array of strings.

# CSP.Image

Provides access to the Image part of the Content Security Policy. This layer contains two fields:

- [Image.Default](#image.default)
- [Image.SRC](#image.src)

# Image.Default

Sets the default values for the policy. This value is an array of strings.

# Image.SRC

Sets the SRC values for the policy. This value is an array of strings.