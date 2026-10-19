# Import using the composite ID format: project_id/feature_name
terraform import unleash_feature_v2.with_env_strategies default/my_nice_feature

# Note: After import, environments are stored sorted alphabetically by name
# and tags are sorted by type then value. Declare blocks in the same order
# in your configuration to avoid ordering-only diffs.
