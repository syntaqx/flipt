package flipt

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.flipt.io/flipt/errors"
)

const maxVariantAttachmentSize = 10000

// Validator validates types
type Validator interface {
	Validate() error
}

// Evaluate

func validateAttachment(attachment string) error {
	if attachment == "" {
		return nil
	}

	bytes := []byte(attachment)
	if !json.Valid(bytes) {
		return errors.InvalidFieldError("attachment", "must be a json string")
	}

	if len(bytes) > maxVariantAttachmentSize {
		return errors.InvalidFieldError("attachment",
			fmt.Sprintf("must be less than %d KB", maxVariantAttachmentSize),
		)
	}
	return nil
}

func (req *EvaluationRequest) Validate() error {
	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	if req.EntityId == "" {
		return errors.EmptyFieldError("entityId")
	}

	return nil
}

var keyRegex = regexp.MustCompile(`^[-_,A-Za-z0-9]+$`)

// Flags

func (req *GetFlagRequest) Validate() error {
	if req.Key == "" {
		return errors.EmptyFieldError("key")
	}

	return nil
}

func (req *ListFlagRequest) Validate() error {
	if req.Limit == 0 && (req.Offset > 0 || req.PageToken != "") {
		return errors.ErrInvalid("limit must be set when offset or pageToken is set")
	}

	return nil
}

func (req *CreateFlagRequest) Validate() error {
	if req.Key == "" {
		return errors.EmptyFieldError("key")
	}

	if !keyRegex.MatchString(req.Key) {
		return errors.InvalidFieldError("key", "contains invalid characters")
	}

	if req.Name == "" {
		return errors.EmptyFieldError("name")
	}

	return nil
}

func (req *UpdateFlagRequest) Validate() error {
	if req.Key == "" {
		return errors.EmptyFieldError("key")
	}

	if req.Name == "" {
		return errors.EmptyFieldError("name")
	}

	return nil
}

func (req *DeleteFlagRequest) Validate() error {
	if req.Key == "" {
		return errors.EmptyFieldError("key")
	}

	return nil
}

func (req *CreateVariantRequest) Validate() error {
	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	if req.Key == "" {
		return errors.EmptyFieldError("key")
	}

	if err := validateAttachment(req.Attachment); err != nil {
		return err
	}

	return nil
}

func (req *UpdateVariantRequest) Validate() error {
	if req.Id == "" {
		return errors.EmptyFieldError("id")
	}

	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	if req.Key == "" {
		return errors.EmptyFieldError("key")
	}

	if err := validateAttachment(req.Attachment); err != nil {
		return err
	}

	return nil
}

func (req *DeleteVariantRequest) Validate() error {
	if req.Id == "" {
		return errors.EmptyFieldError("id")
	}

	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	return nil
}

// Rules

func (req *ListRuleRequest) Validate() error {
	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	if req.Limit == 0 && (req.Offset > 0 || req.PageToken != "") {
		return errors.ErrInvalid("limit must be set when offset or pageToken is set")
	}

	return nil
}

func (req *GetRuleRequest) Validate() error {
	if req.Id == "" {
		return errors.EmptyFieldError("id")
	}

	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	return nil
}

func (req *CreateRuleRequest) Validate() error {
	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	if req.SegmentKey == "" && len(req.SegmentKeys) == 0 {
		return errors.EmptyFieldError("segmentKey or segmentKeys")
	}

	if req.SegmentKey != "" && len(req.SegmentKeys) > 0 {
		return errors.InvalidFieldError("segmentKey or segmentKeys", "only one can be present")
	}

	if req.Rank <= 0 {
		return errors.InvalidFieldError("rank", "must be greater than 0")
	}

	return nil
}

func (req *UpdateRuleRequest) Validate() error {
	if req.Id == "" {
		return errors.EmptyFieldError("id")
	}

	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	if req.SegmentKey == "" && len(req.SegmentKeys) == 0 {
		return errors.EmptyFieldError("segmentKey or segmentKeys")
	}

	if req.SegmentKey != "" && len(req.SegmentKeys) > 0 {
		return errors.InvalidFieldError("segmentKey or segmentKeys", "only one can be present")
	}

	return nil
}

func (req *DeleteRuleRequest) Validate() error {
	if req.Id == "" {
		return errors.EmptyFieldError("id")
	}

	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	return nil
}

func (req *OrderRulesRequest) Validate() error {
	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	if len(req.RuleIds) < 2 {
		return errors.InvalidFieldError("ruleIds", "must contain atleast 2 elements")
	}

	return nil
}

func (req *CreateDistributionRequest) Validate() error {
	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	if req.RuleId == "" {
		return errors.EmptyFieldError("ruleId")
	}

	if req.VariantId == "" {
		return errors.EmptyFieldError("variantId")
	}

	if req.Rollout < 0 {
		return errors.InvalidFieldError("rollout", "must be greater than or equal to '0'")
	}

	if req.Rollout > 100 {
		return errors.InvalidFieldError("rollout", "must be less than or equal to '100'")
	}

	return nil
}

func (req *UpdateDistributionRequest) Validate() error {
	if req.Id == "" {
		return errors.EmptyFieldError("id")
	}

	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	if req.RuleId == "" {
		return errors.EmptyFieldError("ruleId")
	}

	if req.VariantId == "" {
		return errors.EmptyFieldError("variantId")
	}

	if req.Rollout < 0 {
		return errors.InvalidFieldError("rollout", "must be greater than or equal to '0'")
	}

	if req.Rollout > 100 {
		return errors.InvalidFieldError("rollout", "must be less than or equal to '100'")
	}

	return nil
}

func (req *DeleteDistributionRequest) Validate() error {
	if req.Id == "" {
		return errors.EmptyFieldError("id")
	}

	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	if req.RuleId == "" {
		return errors.EmptyFieldError("ruleId")
	}

	if req.VariantId == "" {
		return errors.EmptyFieldError("variantId")
	}

	return nil
}

// Segments

func (req *GetSegmentRequest) Validate() error {
	if req.Key == "" {
		return errors.EmptyFieldError("key")
	}

	return nil
}

func (req *ListSegmentRequest) Validate() error {
	if req.Limit == 0 && (req.Offset > 0 || req.PageToken != "") {
		return errors.ErrInvalid("limit must be set when offset or pageToken is set")
	}

	return nil
}

func (req *CreateSegmentRequest) Validate() error {
	if req.Key == "" {
		return errors.EmptyFieldError("key")
	}

	if !keyRegex.MatchString(req.Key) {
		return errors.InvalidFieldError("key", "contains invalid characters")
	}

	if req.Name == "" {
		return errors.EmptyFieldError("name")
	}

	return nil
}

func (req *UpdateSegmentRequest) Validate() error {
	if req.Key == "" {
		return errors.EmptyFieldError("key")
	}

	if req.Name == "" {
		return errors.EmptyFieldError("name")
	}

	return nil
}

func (req *DeleteSegmentRequest) Validate() error {
	if req.Key == "" {
		return errors.EmptyFieldError("key")
	}

	return nil
}

func (req *CreateConstraintRequest) Validate() error {
	if req.SegmentKey == "" {
		return errors.EmptyFieldError("segmentKey")
	}

	if req.Property == "" {
		return errors.EmptyFieldError("property")
	}

	if req.Operator == "" {
		return errors.EmptyFieldError("operator")
	}

	operator := strings.ToLower(req.Operator)
	// validate operator works for this constraint type
	switch req.Type {
	case ComparisonType_STRING_COMPARISON_TYPE:
		if _, ok := StringOperators[operator]; !ok {
			return errors.ErrInvalidf("constraint operator %q is not valid for type string", req.Operator)
		}
	case ComparisonType_NUMBER_COMPARISON_TYPE:
		if _, ok := NumberOperators[operator]; !ok {
			return errors.ErrInvalidf("constraint operator %q is not valid for type number", req.Operator)
		}
	case ComparisonType_BOOLEAN_COMPARISON_TYPE:
		if _, ok := BooleanOperators[operator]; !ok {
			return errors.ErrInvalidf("constraint operator %q is not valid for type boolean", req.Operator)
		}
	case ComparisonType_DATETIME_COMPARISON_TYPE:
		if _, ok := NumberOperators[operator]; !ok {
			return errors.ErrInvalidf("constraint operator %q is not valid for type datetime", req.Operator)
		}
	default:
		return errors.ErrInvalidf("invalid constraint type: %q", req.Type.String())
	}

	if req.Value == "" {
		// check if value is required
		if _, ok := NoValueOperators[operator]; !ok {
			return errors.EmptyFieldError("value")
		}
	} else if req.Type == ComparisonType_DATETIME_COMPARISON_TYPE {
		// we know that a value is set and that the type is datetime
		// so validate that the value is a valid datetime
		// also convert it to UTC before we save
		// TODO: don't love that we are doing this here
		v, err := tryParseDateTime(req.Value)
		if err != nil {
			return err
		}
		req.Value = v
	}

	return nil
}

func (req *UpdateConstraintRequest) Validate() error {
	if req.Id == "" {
		return errors.EmptyFieldError("id")
	}

	if req.SegmentKey == "" {
		return errors.EmptyFieldError("segmentKey")
	}

	if req.Property == "" {
		return errors.EmptyFieldError("property")
	}

	if req.Operator == "" {
		return errors.EmptyFieldError("operator")
	}

	operator := strings.ToLower(req.Operator)
	// validate operator works for this constraint type
	switch req.Type {
	case ComparisonType_STRING_COMPARISON_TYPE:
		if _, ok := StringOperators[operator]; !ok {
			return errors.ErrInvalidf("constraint operator %q is not valid for type string", req.Operator)
		}
	case ComparisonType_NUMBER_COMPARISON_TYPE:
		if _, ok := NumberOperators[operator]; !ok {
			return errors.ErrInvalidf("constraint operator %q is not valid for type number", req.Operator)
		}
	case ComparisonType_BOOLEAN_COMPARISON_TYPE:
		if _, ok := BooleanOperators[operator]; !ok {
			return errors.ErrInvalidf("constraint operator %q is not valid for type boolean", req.Operator)
		}
	case ComparisonType_DATETIME_COMPARISON_TYPE:
		if _, ok := NumberOperators[operator]; !ok {
			return errors.ErrInvalidf("constraint operator %q is not valid for type datetime", req.Operator)
		}
	default:
		return errors.ErrInvalidf("invalid constraint type: %q", req.Type.String())
	}

	if req.Value == "" {
		// check if value is required
		if _, ok := NoValueOperators[operator]; !ok {
			return errors.EmptyFieldError("value")
		}
	} else if req.Type == ComparisonType_DATETIME_COMPARISON_TYPE {
		// we know that a value is set and that the type is datetime
		// so validate that the value is a valid datetime
		// also convert it to UTC before we save
		// TODO: don't love that we are doing this here
		v, err := tryParseDateTime(req.Value)
		if err != nil {
			return err
		}
		req.Value = v
	}

	return nil
}

func (req *DeleteConstraintRequest) Validate() error {
	if req.Id == "" {
		return errors.EmptyFieldError("id")
	}

	if req.SegmentKey == "" {
		return errors.EmptyFieldError("segmentKey")
	}

	return nil
}

// Namespaces
func (req *CreateNamespaceRequest) Validate() error {
	if req.Key == "" {
		return errors.EmptyFieldError("key")
	}

	if !keyRegex.MatchString(req.Key) {
		return errors.InvalidFieldError("key", "contains invalid characters")
	}

	if req.Name == "" {
		return errors.EmptyFieldError("name")
	}

	return nil
}

func (req *UpdateNamespaceRequest) Validate() error {
	if req.Key == "" {
		return errors.EmptyFieldError("key")
	}

	if req.Name == "" {
		return errors.EmptyFieldError("name")
	}

	return nil
}

func (req *CreateRolloutRequest) Validate() error {
	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	switch rule := req.Rule.(type) {
	case *CreateRolloutRequest_Threshold:
		if rule.Threshold.Percentage < 0 || rule.Threshold.Percentage > 100.0 {
			return errors.InvalidFieldError("threshold.percentage", "must be within range [0, 100]")
		}
	case *CreateRolloutRequest_Segment:
		if rule.Segment.SegmentKey == "" && len(rule.Segment.SegmentKeys) == 0 {
			return errors.EmptyFieldError("segmentKey or segmentKeys")
		}

		if rule.Segment.SegmentKey != "" && len(rule.Segment.SegmentKeys) > 0 {
			return errors.InvalidFieldError("segmentKey or segmentKeys", "only one can be present")
		}
	}

	return nil
}

func (req *UpdateRolloutRequest) Validate() error {
	if req.Id == "" {
		return errors.EmptyFieldError("id")
	}

	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	switch rule := req.Rule.(type) {
	case *UpdateRolloutRequest_Threshold:
		if rule.Threshold.Percentage < 0 || rule.Threshold.Percentage > 100.0 {
			return errors.InvalidFieldError("threshold.percentage", "must be within range [0, 100]")
		}
	case *UpdateRolloutRequest_Segment:
		if rule.Segment.SegmentKey == "" && len(rule.Segment.SegmentKeys) == 0 {
			return errors.EmptyFieldError("segmentKey or segmentKeys")
		}

		if rule.Segment.SegmentKey != "" && len(rule.Segment.SegmentKeys) > 0 {
			return errors.InvalidFieldError("segmentKey or segmentKeys", "only one can be present")
		}

	}

	return nil
}

func (req *DeleteRolloutRequest) Validate() error {
	if req.Id == "" {
		return errors.EmptyFieldError("id")
	}

	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	return nil
}

func (req *OrderRolloutsRequest) Validate() error {
	if req.FlagKey == "" {
		return errors.EmptyFieldError("flagKey")
	}

	if len(req.RolloutIds) < 2 {
		return errors.InvalidFieldError("rolloutIds", "must contain atleast 2 elements")
	}

	return nil
}

func tryParseDateTime(v string) (string, error) {
	if d, err := time.Parse(time.RFC3339, v); err == nil {
		return d.UTC().Format(time.RFC3339), nil
	}

	if d, err := time.Parse(time.DateOnly, v); err == nil {
		return d.UTC().Format(time.DateOnly), nil
	}

	return "", errors.ErrInvalidf("parsing datetime from %q", v)
}
