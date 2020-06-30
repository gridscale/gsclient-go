• Branch release branch of develop
• Finalise changelog
• Bump version number (if it exists)
• Push release branch (in case we want to do hotfixes later)
• Merge --no-ff release branch onto master ((master) $
 git merge --no-ff release/release-2.1.0
 git push upstream new-release-2.1.01
• create release on master with tag following the naming scheme for the project (I stick to semver)
• Merge --no-ff release branch onto develop ((develop) $ 
git merge --no-ff release/release-2.1.0
git push upstream new-release-2.1.0
• Unlock changelog in develop