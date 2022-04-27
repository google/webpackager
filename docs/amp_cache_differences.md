# Advice for sites with a mix of AMP and non-AMP

The Google SXG Cache requirements differ from the earlier [Google AMP Cache
requirements](https://github.com/ampproject/amppackager/blob/releases/docs/cache_requirements.md)
for signed exchanges.

The original reasoning for the latter was to validate that AMP SXGs were at
parity with the stronger constraints of AMP. Today, if an AMP page is signed,
but it meets the new SXG Cache requirements instead of the AMP Cache
requirements, its SXG is considered invalid. However, the AMP itself may still
be considered valid, and served like typical AMP. Some tools like [AMP
Packager](https://github.com/ampproject/amppackager) and [Cloudflare AMP Real
URL](https://blog.cloudflare.com/announcing-amp-real-url/) are designed to meet
the AMP SXG requirements.

Ideally, it would be possible to sign an AMP page that meets the new SXG Cache
requirements, and have it function both as an SXG on supporting browsers, and
as an AMP page on other browsers. Thus, new, more general tools like [Web
Packager](https://github.com/google/webpackager) and [Cloudflare Automatic
Signed Exchanges](https://blog.cloudflare.com/automatic-signed-exchanges/)
would function on sites with AMP.

In the meantime, web publishers must decide. If their site is primarily AMP and
they wish to fix their URLs, use an AMP tool. If their site is primarily
non-AMP and they wish to speed up page loads with prefetching, use a general
tool.

For sites with a mix of AMP and non-AMP, they could run both, and use URL
patterns to forward requests as appropriate.

Alternatively, to speed up non-AMP page loads but not fix AMP URLs, it may be
acceptable to run a general tool only. The only downside is the presence of
[Search Console
warnings](https://support.google.com/webmasters/answer/7450883#sgx_warning_list)
that indicate that AMP is being treated as typical, instead of as an SXG.

It may be possible to meet the AMP SXG requirements using a general tool,
though this is not well-tested. For instance:

 - Run [AMP
   Optimizer](https://github.com/ampproject/amp-toolbox/tree/master/packages/optimizer)
   on AMP pages in order to meet the transformed requirement.
 - Modify AMP Optimizer's output to add the [`data-sxg-no-header`
   attribute](https://github.com/google/sxg-rs/blob/main/README.md#preload-subresources)
   to any `<link rel=preload>` tags. This will prevent Cloudflare Automatic
   Signed Exchanges from turning them into Link headers that wouldn't meet the
   AMP SXG requirements.
 - Set a `Cache-Control: max-age` of at least `345600` (4 days).
 - Additional changes may be required; please suggest updates to this doc if
   you discover any.
