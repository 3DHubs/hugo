# How to create a new release of 3dhubs-hugo

1. Install [goreleaser](https://goreleaser.com/install/)
2. Clone this repo onto your local machine
3. Add the main repo as a remote. `git remote add upstream https://github.com/gohugoio/hugo`
4. Fetch the upstream code `git fetch upstream`
5. Merge the latest release tag `git merge v55.1.1`
6. Tag the merge commit yourself `git tag v55.1.1-3dhubs`
7. Run goreleaser: `goreleaser`
8. Download the latest versions to your dev machines
9. Update the [Dockerfile](https://github.com/3DHubs/util-hugo/blob/master/Dockerfile) for `util-hugo`.

That should be it 🎉
