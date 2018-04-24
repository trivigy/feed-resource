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
  type: feed-resource
  source:
    url: https://example.com/rss.xml
    insecure: true
```

Behavior
--------

### `check`: Checks for a new version of the feed and versions its hash.

#### Source: `source`

- `url` - (_required_) The URL of the feed to consume. The feed file needs to 
    comply with the [gofeed](https://github.com/mmcdole/gofeed) library standard.
- `insecure` - (_optional_) [false] Skip verification of remote TLS certificates.

### `in`: Downloads the new feed as saves it in a file named after the feed name.

### `out`: Inactive.

## Usage Example
```yaml
jobs:
- name: feed-job
  plan:
  - get: my-blog
    trigger: true
    version: every
  - task: build
    config:
      platform: linux
      image_resource:
        type: docker-image
        source:
          repository: ubuntu
          tag: "16.04"
      inputs:
      - name: my-blog
      run:
        path: /bin/bash
        args:
        - -c
        - |
          cd my-blog
          cat rss.xml
```