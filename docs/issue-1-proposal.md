# Proposal: Tenant-related user/member delete (GDPR)

Goal
------
When a tenant requests deletion of a user/member, all personal data related to that user must be removed or irreversibly anonymized while preserving auditability where legally required.

Inputs
------
- tenant_id (bigint)
- user_id (bigint)
- requestor (actor) info derived from JWT/session

Outputs
-------
- JSON report { status: success|partial|failed, total_affected: int, errors: [] }

Features
-------
1. delete_user_gdpr(p_tenant_id, p_user_id) - SQL function (transact)
2. anonymize_admin_audit_for_user(p_tenant_id, p_user_id) - SQL helper
3. Admin HTTP handler to trigger deletion
4. Acceptance tests (Gherkin)
5. Documentation (UML/text + Issue text)

Gherkin Example
---------------
See tests/features/user_delete.feature

Acceptance Criteria
-------------------
- Tenant admin can trigger deletion
- Personal DB fields (email, firstname, lastname) removed or redacted
- Audit rows retained but actor/diff anonymized
- Operation reports success/failure per sub-step

Legal Notes
-----------
- Audit logs may need to be retained for compliance; anonymization/pseudonymization is the recommended approach instead of deletion in many jurisdictions.
