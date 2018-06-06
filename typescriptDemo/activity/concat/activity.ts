import {Observable} from "rxjs/Observable";
import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IFieldDefinition,
    IActivityContribution,
    ActionResult,
    IActionResult
} from "wi-studio/app/contrib/wi-contrib";

@WiContrib({})
@Injectable()
export class ConcatActivityContributionHandler extends WiServiceHandlerContribution {
    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        if(fieldName === "separator") {
           let list: Array<string> = ["-", "$", "#"];
           return list;
        }
        else if(fieldName === "charType") {
            let list: Array<string> = ["Number", "Letter"];
            return list;
         }
        else if(fieldName === "options") {
            let charTypeFieldDef: IFieldDefinition = context.getField("charType");
            if(charTypeFieldDef.value){
                return Observable.create(observer => {
                if(charTypeFieldDef.value=="Number") {
                    var newArray : string[] = new Array();
                    for(let i = 0; i < 10; i++){
                        newArray.push(String(i))
                    }
                    observer.next(newArray);
                }
                else if(charTypeFieldDef.value=="Letter"){
                    var newArray : string[] = new Array();
                    newArray.push("A")
                    newArray.push("B")
                    newArray.push("C")
                    observer.next(newArray);
                }
                
                
            });
        }
        }
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
       if (fieldName === "separator") {
         let vresult: IValidationResult = ValidationResult.newValidationResult();
         let useSeparatorFieldDef: IFieldDefinition = context.getField("useSeparator"); 
         let separatorFieldDef: IFieldDefinition = context.getField("separator");
         if (useSeparatorFieldDef.value && useSeparatorFieldDef.value === true) {
             if (separatorFieldDef.display && separatorFieldDef.display.visible == false) {
                 vresult.setVisible(true);
             } 
             if (separatorFieldDef.value === null || separatorFieldDef.value === "") {
               vresult.setError("TIBCO-CONCAT-1000","Separator must be configured");
             } 
         } else {
            vresult.setVisible(false);
         }
         return vresult;
       }
      return null; 
    }
}