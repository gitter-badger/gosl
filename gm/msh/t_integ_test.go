// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
)

func TestInteg01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Integ01. integration of scalar function")

	// vertices (diamond shape)
	X := [][]float64{
		{0.0, +0.0},
		{1.0, -1.0},
		{2.0, +0.0},
		{1.0, +1.0},
	}

	// allocate cell integrator with default integration points
	o, err := NewIntegrator(TypeQua4, X, nil, "")
	if err != nil {
		tst.Errorf("%v", err)
		return
	}
	chk.Int(tst, "Nverts", o.Nverts, 4)
	chk.Int(tst, "Ndim", o.Ndim, 2)
	chk.Int(tst, "Npts", o.Npts, 4)

	// integrand function
	fcn := func(x []float64) (f float64, e error) {
		f = x[0]*x[0] + x[1]*x[1]
		return
	}

	// perform integration
	res, err := o.IntegrateSv(fcn)
	if err != nil {
		tst.Errorf("%v", err)
		return
	}
	io.Pforan("1: res = %v\n", res)
	chk.Scalar(tst, "∫(x²+y²)dxdy (default)", 1e-15, res, 8.0/3.0)

	// reset integration points
	err = o.ResetP(nil, "legendre_9")
	if err != nil {
		tst.Errorf("%v", err)
		return
	}

	// perform integration again
	res, err = o.IntegrateSv(fcn)
	if err != nil {
		tst.Errorf("%v", err)
		return
	}
	io.Pforan("2: res = %v\n", res)
	chk.Scalar(tst, "∫(x²+y²)dxdy (legendre 9)", 1e-15, res, 8.0/3.0)

	// reset integration points
	err = o.ResetP(nil, "wilson5corner_5")
	if err != nil {
		tst.Errorf("%v", err)
		return
	}

	// perform integration again
	res, err = o.IntegrateSv(fcn)
	if err != nil {
		tst.Errorf("%v", err)
		return
	}
	io.Pforan("3: res = %v\n", res)
	chk.Scalar(tst, "∫(x²+y²)dxdy (wilson5corner)", 1e-15, res, 8.0/3.0)

	// draw polygon
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.Polyline(X, &plt.A{C: "#f4c392", L: "curve1", NoClip: true})
		for _, x := range o.Xip {
			plt.PlotOne(x[0], x[1], &plt.A{C: "b", M: "o", Ms: 6, NoClip: true})
		}
		plt.Gll("x", "y", nil)
		plt.AxisRange(0, 2, -1, 1)
		plt.Equal()
		plt.HideTRborders()
		plt.Save("/tmp/gosl/gm", "integ01")
	}
}

func TestInteg02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Integ02. integration of scalar function")

	// vertices (trapezium)
	a, b, h := 3.0, 0.5, 1.0
	X := [][]float64{
		{-a / 2.0, -h / 2.0},
		{+a / 2.0, -h / 2.0},
		{+b / 2.0, +h / 2.0},
		{-b / 2.0, +h / 2.0},
	}

	// allocate cell integrator with default integration points
	o, err := NewIntegrator(TypeQua4, X, nil, "legendre_4")
	if err != nil {
		tst.Errorf("%v", err)
		return
	}

	// integrand function for moment of inertia about x-axis: Ix
	ρ := 1.0 // density distribution; could be a function ρ(x,y)
	fcnIx := func(x []float64) (f float64, e error) {
		f = x[1] * x[1] * ρ
		return
	}

	// integrand function for moment of inertia about y-axis: Iy
	fcnIy := func(x []float64) (f float64, e error) {
		f = x[0] * x[0] * ρ
		return
	}

	// integrand function for moment of inertia about the origin: I0
	fcnI0 := func(x []float64) (f float64, e error) {
		f = (x[0]*x[0] + x[1]*x[1]) * ρ
		return
	}

	// analytical solutions
	anaIx := (a + b) * math.Pow(h, 3) / 24.0
	anaIy := h * (math.Pow(a, 4) - math.Pow(b, 4)) / (48.0 * (a - b))
	anaI0 := anaIx + anaIy

	// compute Ix
	Ix, err := o.IntegrateSv(fcnIx)
	if err != nil {
		tst.Errorf("%v", err)
		return
	}
	io.Pforan("Ix = %v\n", Ix)
	chk.Scalar(tst, "Ix", 1e-15, Ix, anaIx)

	// compute Iy
	Iy, err := o.IntegrateSv(fcnIy)
	if err != nil {
		tst.Errorf("%v", err)
		return
	}
	io.Pforan("Iy = %v\n", Iy)
	chk.Scalar(tst, "Iy", 1e-15, Iy, anaIy)

	// compute I0
	I0, err := o.IntegrateSv(fcnI0)
	if err != nil {
		tst.Errorf("%v", err)
		return
	}
	io.Pforan("I0 = %v\n", I0)
	chk.Scalar(tst, "I0", 1e-15, I0, anaI0)

	// draw polygon
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		plt.Polyline(X, &plt.A{C: "#f4c392", L: "curve1", NoClip: true})
		for _, x := range o.Xip {
			plt.PlotOne(x[0], x[1], &plt.A{C: "b", M: "o", Ms: 6, NoClip: true})
		}
		plt.Gll("x", "y", nil)
		plt.AxisRange(0, 2, -1, 1)
		plt.Equal()
		plt.HideTRborders()
		plt.Save("/tmp/gosl/gm", "integ02")
	}
}
