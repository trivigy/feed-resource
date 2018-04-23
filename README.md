# RSS Feed Concourse Resource

Retrieves and parses an RSS/Atom feed from an arbitrary
URL, and splits out each syndicated post into files on-disk.

Resource Type Configuration
---------------------------

```yaml
resource_types:
- name: feed-resource
  type: docker-image
  source:
    repository: syncaide/feed-resource
```

Source Configuration
--------------------

```yaml
resources:
- name: my-blog
  type: rss
  source:
    url: https://example.com/rss.xml
```

The following `source` properties are defined:

- `url` - (_required_) The URL of the RSS feed to consume.
- `skip_tls_verify` - (_optional_) [false] Skip verification of remote TLS
  certificates.