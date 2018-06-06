import {NgModule} from "@angular/core";
import {TimerTriggerService} from "./timer";
import {WiServiceContribution} from "wi-studio/app/contrib/wi-contrib";

@NgModule({
    providers: [
        {
            provide: WiServiceContribution,
            useClass: TimerTriggerService
        }
    ],
})
export default class TimerContribModule {

}
