package planmodifiers

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func RequiresReplaceIfValuesNotNull() planmodifier.Map {
	return requiresReplaceIfValuesNotNullModifier{}
}

type requiresReplaceIfValuesNotNullModifier struct{}

func (r requiresReplaceIfValuesNotNullModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.State.Raw.IsNull() {
		// Terraform State의 Values가 Null인 경우
		// Resource를 생성하는 경우이므로 Replace가 필요하지 않음
		return
	}

	if req.Plan.Raw.IsNull() {
		// Terrafomr Plan에 전달된 Values가 Null인 경우
		// Resource를 삭제하는 경우이므로 Replace가 필요하지 않음
		return
	}

	if req.ConfigValue.Equal(req.StateValue) {
		// Terraform 에서 전달된 Values와 StateValue가 동일한 경우
		// Plan에 변경사항이 없으므로 Replace가 필요하지 않음
		return
	}

	if req.StateValue.IsNull() {
		// Terraform State가 Null인 경우, Terraform 최초 배포시

		// 전달된 Config Value가 하나라도 존재한다면, 상태를 변경한다.
		allNullValues := true
		for _, configValue := range req.ConfigValue.Elements() {
			if !configValue.IsNull() {
				allNullValues = false
			}
		}

		if allNullValues {
			return
		}
	} else {
		allNewNullValues := true

		configMap := req.ConfigValue
		stateMap := req.StateValue

		for configKey, configValues := range configMap.Elements() {
			stateValue, ok := stateMap.Elements()[configKey]

			if !ok && configValues.IsNull() {
				// Terraform State에 해당 Key가 없고, Config에 해당 Key가 Null인 경우, Replace가 필요하지 않음
				continue
			}

			if configValues.Equal(stateValue) {
				// Config와 Terraform State의 값이 동일한 경우, Replace가 필요하지 않음
				continue
			}

			allNewNullValues = false
			break
		}

		for stateKey := range stateMap.Elements() {
			_, ok := configMap.Elements()[stateKey]

			if !ok {
				// Terraform State에만 존재하는 Key가 있는 경우 Replace가 필요함
				allNewNullValues = false
				break
			}
		}

		if allNewNullValues {
			return
		}
	}

	resp.RequiresReplace = true
}

func (r requiresReplaceIfValuesNotNullModifier) Description(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (r requiresReplaceIfValuesNotNullModifier) MarkdownDescription(ctx context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}
