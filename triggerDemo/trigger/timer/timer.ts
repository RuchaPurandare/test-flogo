import {Injectable} from "@angular/core";
import {
    WiContrib,
    WiServiceHandlerContribution
} from "wi-studio/app/contrib/wi-contrib";
import {ITriggerContribution, IFieldDefinition} from "wi-studio/common/models/contrib";
import {IValidationResult, ValidationResult} from "wi-studio/common/models/validation";
import {Observable} from "rxjs/Observable";

@Injectable()
@WiContrib({})

export class TimerTriggerService extends WiServiceHandlerContribution {
    constructor() {
        super();
    }
    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
        return null;
    }

    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "Time Interval" || fieldName === "Interval Unit") {
            let field: IFieldDefinition = context.getField("Repeating");
            if (field.value === false) {
                return ValidationResult.newValidationResult().setVisible(false);
            } else {
                if (fieldName === "Time Interval") {
                  let interval: number = context.getField("Time Interval").value;
                  if (interval <= 0) {
                    return ValidationResult.newValidationResult().setVisible(true).setError("GENERAL-TIMER-1000","Interval value must be a positive number and must be greater than zero");
                  }
                }
                return ValidationResult.newValidationResult().setVisible(true);
            }
        }
        return null;
    }
}
