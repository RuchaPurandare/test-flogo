import { Injectable } from "@angular/core";
import {WiContrib, WiServiceHandlerContribution, CreateFlowActionResult, ICreateFlowActionContext,
    IActionResult,ActionResult} from "wi-studio/app/contrib/wi-contrib";
import {ITriggerContribution, IFieldDefinition} from "wi-studio/common/models/contrib";
import {IValidationResult} from "wi-studio/common/models/validation";
import {Observable} from "rxjs/Observable";
import * as lodash from "lodash";

@Injectable()
@WiContrib({})
export class SampleActivityUIContributionHandler extends WiServiceHandlerContribution {
    constructor() {
        super();
    }
 
    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
        return null;
    }
  
    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
       return null;
    }

    action = (fieldName: string, context: ICreateFlowActionContext): Observable<IActionResult> | IActionResult => {
        let modelService = this.getModelService();
        let result = CreateFlowActionResult.newActionResult();
        /*if (context.handler && context.handler.settings && context.handler.settings.length > 0) {
            let methodValue = <IFieldDefinition>context.getField("Method");
            let output = <IFieldDefinition>context.getField("body");
            // Create multiple flows based on Method field selection on wizard 
            if (methodValue && methodValue.value && methodValue.value.length > 0) {
                for (let i = 0; i < methodValue.value.length; i++) {
                    let trigger = modelService.createTriggerElement("General/tibco-wi-rest");
                    // Set the field values for your trigger 
                    if (trigger && trigger.handler && trigger.handler.settings && trigger.handler.settings.length > 0) {
                        for (let j = 0; j < trigger.handler.settings.length; j++) {
                            if (trigger.handler.settings[j].name === "Method") {
                                trigger.handler.settings[j].value = methodValue.value[i];
                            }
                        }
                    }
                    if (trigger && trigger.outputs && trigger.outputs.length > 0) {
                        for (let j = 0; j < trigger.outputs.length; j++) {
                            if (trigger.outputs[j].name === "body") {
                                trigger.outputs[j].value = output.value;
                            }
                        }
                    }
                     //Add new activity for each flow
                    // Set the field values for your activity
                    let reply = modelService.createFlowElement("General/tibco-wi-reply");
                    if (reply && reply.inputs && reply.inputs.length > 0) {
                        for (let j = 0; j < reply.inputs.length; j++) {
                            if (reply.inputs[j].name === "data") {
                                reply.inputs[j].value = output.value;
                            }
                        }
                    }
                    let flowModel = modelService.createFlow(context.getFlowName() + "_" + methodValue.value[i], context.getFlowDescription());
                    let flow = flowModel.addFlowElement(reply);
                    result = result.addTriggerFlowMapping(lodash.cloneDeep(trigger), lodash.cloneDeep(flow));
                }
            }
        } */
        let actionResult = ActionResult.newActionResult().setSuccess(true).setResult(result);
        return actionResult;
    }
}