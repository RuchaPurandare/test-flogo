import { Injector } from "@angular/core";
import { WiServiceContribution, ITriggerContribution, IValidationResult } from "wi-studio/index";
import { HttpModule, Http, XHRBackend } from "@angular/http";
import { } from "jasmine";
import { TestBed, inject } from "@angular/core/testing";
import { MockBackend } from "@angular/http/testing";
import { TimerTriggerService } from "./timer";

/**
 * TimerTriggerService tests
 */
export let t1 = describe("TimerTriggerService tests", () => {
    beforeEach(() => {
        TestBed.configureTestingModule({
            imports: [HttpModule],
            providers: [
                { provide: WiServiceContribution, useClass: TimerTriggerService },
                { provide: XHRBackend, useClass: MockBackend }
            ]
        });
    });
    afterEach(() => {
        TestBed.resetTestingModule();
    });


    /**
     * Test TimerTriggerService
     */
    describe("TimerTriggerService", () => {
        it("should return timer trigger service", () => {
            inject([Injector, Http], (injector: Injector, http: Http) => {
                let svc = new TimerTriggerService();
                expect(svc !== null).toBeTruthy("Timer Trigger Service not found");
            })();
        });
    });
    /**
     * Test validate for the Time Interval
     */
    describe("getValidationProvider(Time Interval)", () => {
        it("should return the validation provider for Time Interval", () => {
            inject([Injector, Http], (injector: Injector, http: Http) => {
                let svc = new TimerTriggerService();

                describe("IContribValidation(context:false)", () => {
                    const context: any = {
                        endpoint: {
                            settings: [
                                {
                                    name: "Repeating",
                                    type: "boolean",
                                    description: "Indicates this flow should be run only once or multiple times",
                                    value: true,
                                    required: true
                                }
                            ]
                        }
                    };
                    it("should return true", () => {
                        let res: IValidationResult  = <IValidationResult>svc.validate("Time Interval", <ITriggerContribution>context);
                        expect(res.isVisible()).toBeFalsy("Incorrect value returned, expected false");

                    });
                });

                describe("IContribValidation(context:true)", () => {
                    let context: any = {
                        endpoint: {
                            settings: [
                                {
                                    name: "Repeating",
                                    type: "boolean",
                                    description: "Indicates this flow should be run only once or multiple times",
                                    value: false,
                                    required: true
                                }
                            ]
                        }
                    };
                    it("should return false", () => {
                        let res: IValidationResult  = <IValidationResult>svc.validate("Time Interval", <ITriggerContribution>context);
                        expect(res.isVisible()).toBeFalsy("Incorrect value returned, expected false");
                    });
                });
            })();
        });
    });
    /**
     * Test validate for the Interval Unit
     */
    describe("getValidationProvider(Interval Unit)", () => {
        it("should return a validation provider for Interval Unit", () => {
            inject([Injector, Http], (injector: Injector, http: Http) => {
                let svc = new TimerTriggerService();
                describe("IContribValidation(context:false)", () => {
                    const context: any = {
                        endpoint: {
                            settings: [
                                {
                                    name: "Repeating",
                                    type: "boolean",
                                    description: "Indicates this flow should be run only once or multiple times",
                                    value: true,
                                    required: true
                                }
                            ]
                        }
                    };
                    it("should return true", () => {
                        let res: IValidationResult  = <IValidationResult>svc.validate("Interval Unit", <ITriggerContribution>context);
                        expect(res.isVisible()).toBeFalsy("Incorrect value returned, expected false");
                    });
                });

                describe("IContribValidation(context:true)", () => {
                    let context: any = {
                        endpoint: {
                            settings: [
                                {
                                    name: "Repeating",
                                    type: "boolean",
                                    description: "Indicates this flow should be run only once or multiple times",
                                    value: false,
                                    required: true
                                }
                            ]
                        }
                    };
                    it("should return false", () => {
                        let res: IValidationResult  = <IValidationResult>svc.validate("Interval Unit", <ITriggerContribution>context);
                        expect(res.isVisible()).toBeFalsy("Incorrect value returned, expected false");
                    });
                });
            })();
        });
    });

});
