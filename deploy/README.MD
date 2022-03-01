# GitFlow Actions
Pull Request

    - First time PR creation runs pull_request_builder.yaml
    - Pushing a commit to an existing PR runs pull_request_builder.yaml
    - Merging an existing PR to develop/master runs continuous_deployment.yaml
    - Direct commit push to develop/master runs continuous_deployment.yaml

Release/Hotfix

    - First time release/hotfix creation against develop/master runs pull_request_builder.yaml
    - Pushing a commit to an existing release/hotfix runs pull_request_builder.yaml
    - Creating a new PR against the release/hotfix runs pull_request_builder.yaml
    - Pushing a commit to an existing PR that tracks a release/hotfix runs pull_request_builder.yaml
    - Merging an existing PR to a release/hotfix runs pull_request_builder.yaml
    - Direct commit push to release/hotfix runs pull_request_builder.yaml
    - Merging an existing release/hotfix to develop/master runs continuous_deployment.yaml
    - Merging an existing release/hotfix to develop/master in CLI and pushing develop/master runs continuous_deployment.yaml

Tag

    - Manually creating a tag does nt trigger any pipeline
    - Pushing a tag from local does not trigger any pipeline
    
Manual Deployment

    - Manually deploying a specific git ref (branch/tag/sha) to a specific environment

Others

    - Closed PRs will not trigger any pipeline
    - Rollback is same as manual deployment so deploy a previous tag

