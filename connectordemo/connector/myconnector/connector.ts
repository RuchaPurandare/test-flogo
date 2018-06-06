/**
 * Imports
 */
import {Inject, Injectable, Injector} from "@angular/core";
import {Http} from "@angular/http";
import { WiProxyCORSUtils, WiContrib, WiContributionUtils, WiServiceHandlerContribution, AUTHENTICATION_TYPE } from "wi-studio/app/contrib/wi-contrib";
import { IConnectorContribution, IFieldDefinition, IActionResult, ActionResult, HTTP_METHOD } from "wi-studio/common/models/contrib";
import { Observable } from "rxjs/Observable";
import { IValidationResult, ValidationResult, ValidationError } from "wi-studio/common/models/validation";

/**
 * Main
 */
@WiContrib({})
@Injectable()
export class MyConnectorContribution extends WiServiceHandlerContribution {
    private category: string;
    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
        this.category = "Jira";
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
        return null;
    }
    
    validate = (name: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
        if (name === "Connect") {
            let name : IFieldDefinition;
            let clientId : IFieldDefinition;
            let userName: IFieldDefinition;
            let password: IFieldDefinition;

            for (let configuration of context.settings) {
                if (configuration.name === "name") {
                    name = configuration
                } else if (configuration.name === "clientId") {
                    clientId = configuration
                } else if (configuration.name === "userName"){
                    userName = configuration
                } else if(configuration.name === "password"){
                    password = configuration
                }
            }
            if (name.value && clientId.value && userName.value && password.value) {
                return ValidationResult.newValidationResult().setReadOnly(false)
            } else {
                return ValidationResult.newValidationResult().setReadOnly(true)
            }
        }
        return null;
    }

    action = (actionName: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {
        if (actionName == "Connect") {
            return Observable.create(observer => {

                let currentName : string;

                for (let i = 0; i < context.settings.length; i++) {
                    if (context.settings[i].name === "name") {
                        currentName = context.settings[i].value;
                    }
                }

                let duplicate = false;
               
                WiContributionUtils.getConnections(this.http, this.category).subscribe((conns: IConnectorContribution[]) => {
                    for (let conn of conns) {
                        for (let i = 0; i < conn.settings.length; i++) {
                            if (conn.settings[i].name === "name") {
                                let oldName = conn.settings[i].value;
                                if (oldName === currentName && (WiContributionUtils.getUniqueId(conn) !== WiContributionUtils.getUniqueId(context))) {
                                    duplicate = true;
                                    break;
                                }
                            }
                        }
                    }
                
                    if (duplicate) {
                        observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("MY-CONNECTOR-001", "Connection name already exists")));
                    } else {
                        let clientId = "", userName = "", password = "";
                        for (let configuration of context.settings) {
                            if (configuration.name === "clientId") {
                                clientId = configuration.value
                            } else if (configuration.name === "userName") {
                                userName = configuration.value
                            } else if (configuration.name === "password") {
                                password = configuration.value
                            }
                        }

                        let myURL = "http://demo0279167.mockable.io/testConnection"
            
                        WiProxyCORSUtils.createRequest(this.http, myURL)
                            .addMethod(HTTP_METHOD.POST)
                            .addHeader("Content-Type", "application/json")
                            .addHeader("Authorization", "Basic " + btoa(userName + ":" + password))
                            .send().subscribe(resp => {                               
                                
                                let actionResult = {
                                    context: context,
                                    authType: AUTHENTICATION_TYPE.BASIC,
                                    authData: {}
                                }
                                observer.next(ActionResult.newActionResult().setSuccess(true).setResult(actionResult));
                                console.log("Response is -->",resp)
                            },
                            error => {
                                observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("MY-CONNECTOR-002", "Failed to create connection. Check your configuration.")));
                            }
                        );
                    }
                });
            });
        }
    return null;
    }
}