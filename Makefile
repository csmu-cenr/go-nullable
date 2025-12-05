.PHONY: count push tags version

count:
	./count-git-commits-across-all-branches.sh
	
default: push

push:
	@echo "Pushing to repository using 'git push'"
	@git push
	@echo "Pushing tags to repository using 'git push origin --tags'"
	@git push origin --tags
	
version:
	@echo "Showing versions using 'git tag -n1'"
	@git tag -n1
	
tags:
	git tag -n1 | sort -rh | head -n 5
	
	
