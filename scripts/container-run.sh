#!/bin/sh

aqua policy allow /app/server/aqua-policy.yaml \
&& aqua -c /app/server/aqua.yaml i \
&& mage -d /app/server dev \
& sleep infinity
