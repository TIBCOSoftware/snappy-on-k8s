/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	snappydatav1beta1 "snappydata-operator/api/v1beta1"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// SmallClusterConfiguration contains predefined configuration for "Small" cluster size
var SmallClusterConfiguration = snappydatav1beta1.ClusterConfiguration{
	// Leads
	Leads: snappydatav1beta1.MemberConfiguration{
		Count: 1,
		Conf:  "-heap-size=3g",
		Resources: v1.ResourceRequirements{
			// Limits: apiv1.ResourceList{
			// 	"cpu":    resource.MustParse(cpuLimit),
			// 	"memory": resource.MustParse(memLimit),
			// },
			Requests: v1.ResourceList{
				// "cpu":    resource.MustParse(cpuReq),
				"memory": resource.MustParse("1024Mi"),
			},
		},
		VolumeClaimTemplates: []v1.PersistentVolumeClaim{
			{
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{
						"ReadWriteOnce",
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							"storage": resource.MustParse("10G"),
						},
					},
				},
			},
		},
	},

	// Locators
	Locators: snappydatav1beta1.MemberConfiguration{
		Count: 1,
		Conf:  "-heap-size=3g",
		// Resources: v1.ResourceRequirements{
		// 	// Limits: apiv1.ResourceList{
		// 	// 	"cpu":    resource.MustParse(cpuLimit),
		// 	// 	"memory": resource.MustParse(memLimit),
		// 	// },
		// 	Requests: v1.ResourceList{
		// 		// "cpu":    resource.MustParse(cpuReq),
		// 		"memory": resource.MustParse("1024Mi"),
		// 	},
		// },
		VolumeClaimTemplates: []v1.PersistentVolumeClaim{
			{
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{
						"ReadWriteOnce",
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							"storage": resource.MustParse("10G"),
						},
					},
				},
			},
		},
	},

	// Servers
	Servers: snappydatav1beta1.MemberConfiguration{
		Count: 1,
		Conf:  "-heap-size=3g",
		// Resources: v1.ResourceRequirements{
		// 	// Limits: apiv1.ResourceList{
		// 	// 	"cpu":    resource.MustParse(cpuLimit),
		// 	// 	"memory": resource.MustParse(memLimit),
		// 	// },
		// 	Requests: v1.ResourceList{
		// 		// "cpu":    resource.MustParse(cpuReq),
		// 		"memory": resource.MustParse("1024Mi"),
		// 	},
		// },
		VolumeClaimTemplates: []v1.PersistentVolumeClaim{
			{
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{
						"ReadWriteOnce",
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							"storage": resource.MustParse("10G"),
						},
					},
				},
			},
		},
	},
}

// MediumClusterConfiguration contains predefined configuration for "Medium" cluster size
var MediumClusterConfiguration = snappydatav1beta1.ClusterConfiguration{
	// Leads
	Leads: snappydatav1beta1.MemberConfiguration{
		Count: 1,
		Conf:  "-heap-size=3g",
		Resources: v1.ResourceRequirements{
			// Limits: apiv1.ResourceList{
			// 	"cpu":    resource.MustParse(cpuLimit),
			// 	"memory": resource.MustParse(memLimit),
			// },
			Requests: v1.ResourceList{
				// "cpu":    resource.MustParse(cpuReq),
				"memory": resource.MustParse("1024Mi"),
			},
		},
		VolumeClaimTemplates: []v1.PersistentVolumeClaim{
			{
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{
						"ReadWriteOnce",
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							"storage": resource.MustParse("10G"),
						},
					},
				},
			},
		},
	},

	// Locators
	Locators: snappydatav1beta1.MemberConfiguration{
		Count: 1,
		Conf:  "-heap-size=3g",
		// Resources: v1.ResourceRequirements{
		// 	// Limits: apiv1.ResourceList{
		// 	// 	"cpu":    resource.MustParse(cpuLimit),
		// 	// 	"memory": resource.MustParse(memLimit),
		// 	// },
		// 	Requests: v1.ResourceList{
		// 		// "cpu":    resource.MustParse(cpuReq),
		// 		"memory": resource.MustParse("1024Mi"),
		// 	},
		// },
		VolumeClaimTemplates: []v1.PersistentVolumeClaim{
			{
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{
						"ReadWriteOnce",
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							"storage": resource.MustParse("10G"),
						},
					},
				},
			},
		},
	},

	// Servers
	Servers: snappydatav1beta1.MemberConfiguration{
		Count: 2,
		Conf:  "-heap-size=3g",
		// Resources: v1.ResourceRequirements{
		// 	// Limits: apiv1.ResourceList{
		// 	// 	"cpu":    resource.MustParse(cpuLimit),
		// 	// 	"memory": resource.MustParse(memLimit),
		// 	// },
		// 	Requests: v1.ResourceList{
		// 		// "cpu":    resource.MustParse(cpuReq),
		// 		"memory": resource.MustParse("1024Mi"),
		// 	},
		// },
		VolumeClaimTemplates: []v1.PersistentVolumeClaim{
			{
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{
						"ReadWriteOnce",
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							"storage": resource.MustParse("10G"),
						},
					},
				},
			},
		},
	},
}

// LargeClusterConfiguration contains predefined configuration for "Large" cluster size
var LargeClusterConfiguration = snappydatav1beta1.ClusterConfiguration{
	// Leads
	Leads: snappydatav1beta1.MemberConfiguration{
		Count: 1,
		Conf:  "-heap-size=3g",
		Resources: v1.ResourceRequirements{
			// Limits: apiv1.ResourceList{
			// 	"cpu":    resource.MustParse(cpuLimit),
			// 	"memory": resource.MustParse(memLimit),
			// },
			Requests: v1.ResourceList{
				// "cpu":    resource.MustParse(cpuReq),
				"memory": resource.MustParse("1024Mi"),
			},
		},
		VolumeClaimTemplates: []v1.PersistentVolumeClaim{
			{
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{
						"ReadWriteOnce",
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							"storage": resource.MustParse("10G"),
						},
					},
				},
			},
		},
	},

	// Locators
	Locators: snappydatav1beta1.MemberConfiguration{
		Count: 1,
		Conf:  "-heap-size=3g",
		// Resources: v1.ResourceRequirements{
		// 	// Limits: apiv1.ResourceList{
		// 	// 	"cpu":    resource.MustParse(cpuLimit),
		// 	// 	"memory": resource.MustParse(memLimit),
		// 	// },
		// 	Requests: v1.ResourceList{
		// 		// "cpu":    resource.MustParse(cpuReq),
		// 		"memory": resource.MustParse("1024Mi"),
		// 	},
		// },
		VolumeClaimTemplates: []v1.PersistentVolumeClaim{
			{
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{
						"ReadWriteOnce",
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							"storage": resource.MustParse("10G"),
						},
					},
				},
			},
		},
	},

	// Servers
	Servers: snappydatav1beta1.MemberConfiguration{
		Count: 4,
		Conf:  "-heap-size=3g",
		// Resources: v1.ResourceRequirements{
		// 	// Limits: apiv1.ResourceList{
		// 	// 	"cpu":    resource.MustParse(cpuLimit),
		// 	// 	"memory": resource.MustParse(memLimit),
		// 	// },
		// 	Requests: v1.ResourceList{
		// 		// "cpu":    resource.MustParse(cpuReq),
		// 		"memory": resource.MustParse("1024Mi"),
		// 	},
		// },
		VolumeClaimTemplates: []v1.PersistentVolumeClaim{
			{
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{
						"ReadWriteOnce",
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							"storage": resource.MustParse("10G"),
						},
					},
				},
			},
		},
	},
}
