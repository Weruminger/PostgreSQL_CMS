Feature: Tenant deletes user and personal data is removed or anonymized
  As a tenant admin
  I want to delete a user and ensure personal data is removed or anonymized
  So that we comply with GDPR

  Scenario: Tenant admin deletes a user and all personal data removed
    Given tenant 1 exists and admin is authenticated
    And a user with id 123 exists with personal data in entries, media_assets and sessions
    When the admin requests POST /admin/tenants/1/users/123/delete
    Then the API returns status 200 and a JSON report with "status":"success"
    And the user row with id 123 does not exist (if users table exists)

  Scenario: Audit entries are anonymized, not deleted
    Given audit entries exist referencing actor "user:123" and diffs containing personal fields
    When user 123 is deleted
    Then audit entries for actions related to user 123 remain
    And actor is replaced with "anon:user:" pattern and personal fields removed from diff
