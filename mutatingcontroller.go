package main

import (
	"encoding/json"
	"errors"
	"gomodules.xyz/jsonpatch/v3"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"log"
)

/**
This code will add a supplied toleration to the pod. This would be useful if you want to add a toleration to all the pods in a namespace.


*/

func MutateCustomAnnotation(admissionRequest *v1beta1.AdmissionRequest ) (*v1beta1.AdmissionResponse, error){

	// Parse the Pod object.
	raw := admissionRequest.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		return nil, errors.New("unable to parse pod")
	}

	// add toleration
	newToleration := corev1.Toleration{Key: "ml-pod", Operator: corev1.TolerationOpEqual, Value: "true", Effect: corev1.TaintEffectNoSchedule}
	pod.Spec.Tolerations = append(pod.Spec.Tolerations, newToleration)

	podWithToleration, err := json.Marshal(pod)
	if err != nil {
		log.Print("error in marshalling the pod object passed ....", err)
		return nil, err
	}


	patchedOperation , err := jsonpatch.CreatePatch ( raw , podWithToleration )
	if err != nil {
		log.Print("error in creatin the patch ....", err)
		return nil, err
	}

	//convert patch into bytes
	patchBytes, err := json.Marshal(patchedOperation)
	if err != nil {
		log.Print("error marshalling the modified pod....", err)
		return nil, errors.New("unable to parse the patch")
	}

	//create the response with patch bytes
	var admissionResponse *v1beta1.AdmissionResponse
	admissionResponse = &v1beta1.AdmissionResponse {
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}

	//return the resonse
	return admissionResponse, nil

}
