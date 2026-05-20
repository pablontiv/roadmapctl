package diagnostics

const (
	DiagnosticSingleFileFallback = "RMC_STRUCTURE_SINGLE_FILE_FALLBACK"
	DiagnosticRootlineMissing    = "RMC_ENV_ROOTLINE_MISSING"
	DiagnosticInvalidBlockedBy   = "RMC_GRAPH_INVALID_BLOCKED_BY"
	DiagnosticConfigMissing      = "RMC_CONFIG_MISSING"
)

const (
	DiagnosticLintTaskTableMissing            = "RMC_LINT_TASK_TABLE_MISSING"
	DiagnosticLintTaskTableMissingRow         = "RMC_LINT_TASK_TABLE_MISSING_ROW"
	DiagnosticLintTaskTableStaleRow           = "RMC_LINT_TASK_TABLE_STALE_ROW"
	DiagnosticLintTaskTableInvalidLink        = "RMC_LINT_TASK_TABLE_INVALID_LINK"
	DiagnosticLintTaskSectionMissing          = "RMC_LINT_TASK_SECTION_MISSING"
	DiagnosticLintAcceptanceCriteriaMissing   = "RMC_LINT_ACCEPTANCE_CRITERIA_MISSING"
	DiagnosticLintSourceOfTruthEmpty          = "RMC_LINT_SOURCE_OF_TRUTH_EMPTY"
	DiagnosticLintFilenameCaseCollision       = "RMC_LINT_FILENAME_CASE_COLLISION"
	DiagnosticLintFilenameReserved            = "RMC_LINT_FILENAME_RESERVED"
	DiagnosticLintSchemaFieldMissing          = "RMC_LINT_SCHEMA_FIELD_MISSING"
	DiagnosticLintSchemaLinkMissing           = "RMC_LINT_SCHEMA_LINK_MISSING"
	DiagnosticLintSchemaOutcomeEstadoRequired = "RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED"
	DiagnosticLintSchemaOutcomeEstadoNonEmpty = "RMC_LINT_SCHEMA_OUTCOME_ESTADO_NON_EMPTY"
)

const (
	DiagnosticTransitionTaskNotFound      = "RMC_TRANSITION_TASK_NOT_FOUND"
	DiagnosticTransitionStatusUnknown     = "RMC_TRANSITION_STATUS_UNKNOWN"
	DiagnosticTransitionDependencyBlocked = "RMC_TRANSITION_DEPENDENCY_BLOCKED"
	DiagnosticTransitionRoleMissing       = "RMC_TRANSITION_ROLE_MISSING"
	DiagnosticTransitionNotActive         = "RMC_TRANSITION_NOT_ACTIVE"
	DiagnosticTransitionAlreadyDone       = "RMC_TRANSITION_ALREADY_DONE"
	DiagnosticTransitionApplyFailed       = "RMC_TRANSITION_APPLY_FAILED"
)

const (
	DiagnosticBootstrapRepairUnsupportedStem = "RMC_BOOTSTRAP_REPAIR_UNSUPPORTED_STEM"
)

const (
	DiagnosticMaterializeInputVersionUnsupported   = "RMC_MATERIALIZE_INPUT_VERSION_UNSUPPORTED"
	DiagnosticMaterializeInputKindInvalid          = "RMC_MATERIALIZE_INPUT_KIND_INVALID"
	DiagnosticMaterializeInputEmpty                = "RMC_MATERIALIZE_INPUT_EMPTY"
	DiagnosticMaterializeInputFieldMissing         = "RMC_MATERIALIZE_INPUT_FIELD_MISSING"
	DiagnosticMaterializeInputSlugInvalid          = "RMC_MATERIALIZE_INPUT_SLUG_INVALID"
	DiagnosticMaterializeInputDependencyInvalid    = "RMC_MATERIALIZE_INPUT_DEPENDENCY_INVALID"
	DiagnosticMaterializeInputDependencyUnresolved = "RMC_MATERIALIZE_INPUT_DEPENDENCY_UNRESOLVED"
	DiagnosticMaterializePlanConflict              = "RMC_MATERIALIZE_PLAN_CONFLICT"
)
