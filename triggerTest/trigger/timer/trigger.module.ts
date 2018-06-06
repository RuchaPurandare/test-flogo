import { NgModule } from "@angular/core";
import { SampleActivityUIContributionHandler} from "./trigger";
import { WiServiceContribution } from "wi-studio/app/contrib/wi-contrib";
 
 
@NgModule({
  providers: [
    {
       provide: WiServiceContribution,
       useClass: SampleActivityUIContributionHandler
     }
  ]
})
 
export default class SampleActivityModule {
 
}