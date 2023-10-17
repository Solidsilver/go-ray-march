// Copyright (c) 2022, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "goki.dev/mat32/v2"

//gosl: hlsl basic
// #include "fastexp.hlsl"
//gosl: end basic

//gosl: start basic

// DataStruct has the test data
type DataStruct struct {

	// raw value
	Raw float32

	// integrated value
	Integ float32

	// exp of integ
	Exp float32

	// must pad to multiple of 4 floats for arrays
	Pad2 float32
}

// ParamStruct has the test params
type ParamStruct struct {

	// rate constant in msec
	Tau float32

	// 1/Tau
	Dt float32

	pad, pad1 float32
}

// IntegFmRaw computes integrated value from current raw value
func (ps *ParamStruct) IntegFmRaw(ds *DataStruct) {
	ds.Integ += ps.Dt * (ds.Raw - ds.Integ)
	ds.Exp = mat32.FastExp(-ds.Integ)
}

//gosl: end basic

// note: only core compute code needs to be in shader -- all init is done CPU-side

func (ps *ParamStruct) Defaults() {
	ps.Tau = 2
	ps.Update()
}

func (ps *ParamStruct) Update() {
	ps.Dt = 1.0 / ps.Tau
}

//gosl: hlsl basic
/*
// // note: double-commented lines required here -- binding is var, set
uniform ParamStruct Params;
[[vk::binding(0, 1)]] RWStructuredBuffer<DataStruct> Data;

[numthreads(64, 1, 1)]

void main(uint3 idx : SV_DispatchThreadID) {
    Params.IntegFmRaw(Data[idx.x]);
}
*/
//gosl: end basic