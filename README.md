# logparser
Simple tool to get statistics when parsing logs

## Format

```
X.Y.Z.W - - [22/May/2018:06:42:27 +0000] "GET /robots.txt HTTP/2.0" 200 3687
```
## To-do

- [ ] Total number of visitors, start date, end date
- [ ] Average requests per minute
- [ ] Top requests per minute (rpm)
- [ ] Check for spammy IPs (IPs sorted by highest rpm)
- [ ] Total bandwidth, bandwidth per page
- [ ] Is spider (does it check for robots.txt)?
- [ ] FAST?
