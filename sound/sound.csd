<CsoundSynthesizer>
<CsOptions>
;Halleluia
;normalize amplitude between yesamp and similarity

;envelope speed (pulsing) switch pulse env with straight line @ 1
;similarity also has filter sweep
;similarity has fm > 0 and filter sweep
</CsOptions>
<CsInstruments>

sr = 44100
ksmps = 100 ;samps / k-period
nchnls = 1
0dbfs = 1
gihandle OSCinit 7770

instr 1
	kf1 init 0
	kf2 init 0
	kf3 init 0
	kf4 init 0
	kf5 init 0
	kf6 init 0
	kf7 init 0
	kf8 init 0
	kCount init 0
	konoff init 0
	ksim init 0
	
	ibasefreq = p4
	iscale = p5
	isimscale = p6
	ifn = p7
	idelbase = p8
	
	icarifn = p9
	imodfreq = p10
	imodamp = p11
	imod_ifn = p12

	kpulse_env_amp = p13
	ipulse_shape = p14
	ilowpulse = p15
	ihighpulse = p16
	isimthresh = p17	

 	kk OSClisten gihandle, "/eeg", "iiiiiiiii", kCount, kf1, kf2, kf3, kf4, kf5, kf6, kf7, kf8
	kk OSClisten gihandle, "/onoff", "i", konoff
	kk OSClisten gihandle, "/similarity", "f", ksim

	printks "%i %i %i %i %i %i %i %i %i\\n", 0, kCount, kf1, kf2, kf3, kf4, kf5, kf6, kf7, kf8	
	printks "toggle is %i\\n", 0, konoff	
	printks "similarity is %f\\n", 0, ksim

if (konoff == 0) then 
	kgoto noamp
	
elseif (ksim  >= isimthresh ) && (konoff == 1) then
	kgoto  similarity ; goto similarity and PULSE
	
elseif (ksim < isimthresh) && (konoff == 1) then
	kgoto yesamp
	
endif

yesamp:
 	;USING DATA FOR AMPLITUDE
 	;currentkf / maxkf
 	
 	asig0 oscili .1 , ibasefreq, 1
 	asig1 oscili (kf1 / 18000000) * iscale, ibasefreq, ifn
 	asig2 oscili (kf2 / 18000000) * iscale, ibasefreq * 2, ifn
 	asig3 oscili (kf3 / 18000000) * iscale, ibasefreq * 3, ifn
 	asig4 oscili (kf4 / 18000000) * iscale, ibasefreq * 4, ifn
  	asig5 oscili (kf5 / 18000000) * iscale, ibasefreq * 5, ifn
 	asig6 oscili (kf6 / 18000000) * iscale, ibasefreq * 6, ifn
 	asig7 oscili (kf7 / 18000000) * iscale, ibasefreq * 7, ifn
 	asig8 oscili (kf8 / 18000000) * iscale, ibasefreq * 8, ifn


 	aoutsig = asig0 + asig1 + asig2 + asig3 + asig4 + asig5 + asig6 + asig7 + asig8 
 	aoutsig = aoutsig * iscale
 	acomp oscili 1, 400, 1
 	abal balance aoutsig, acomp

	afiltsig2 butlp abal, 1000 	
	afiltsig butlp afiltsig2,  1000
	
	adelay1 delay afiltsig, idelbase
	adelay2 delay afiltsig, idelbase * .2
	adelay3 delay afiltsig, idelbase * .3
	adelay4 delay afiltsig, idelbase * .4
	adelay5 delay afiltsig, idelbase * .5	
	adelay6 delay afiltsig, idelbase * .6	
	adelay7 delay afiltsig, idelbase * .7	
	adelay8 delay afiltsig, idelbase * .8	
	
	adelays = adelay1 + adelay2 + adelay3 + adelay4 + adelay5 + adelay6 + adelay7 + adelay8 

 	abal2 balance adelays, acomp
	aout = abal2 
	kgoto output 

similarity:
	;printks "SIMILARITY ** SECTION\\n", 0, ksim
	;scale the incoming ksim to min and max pulses per second
 	kpulse_freq scale ksim, ihighpulse, ilowpulse
 	
	kpulse_env oscili kpulse_env_amp, kpulse_freq, ipulse_shape
	amod oscili imodamp, imodfreq, imod_ifn
	amod2 = amod  * kpulse_env

	acar oscili isimscale * kpulse_env,  ibasefreq + amod2, icarifn
	
	aout = acar 
	kgoto output
	
output:
	out aout


noamp:
	endin



</CsInstruments>
<CsScore>
;f1 0 16384 10 1 0 1 .8 .6 .3 .1
;f1 0 16384 10 1 0 .6 0 .1
f1 0 16384 10 1
f2 0 16384 10 1 0 1 0 .5 0 .2
f4 0 16384 10 1 0 1
f3 0 1024 9 0.5 1 0 ; half sine
f5 0 1025 7  0.01    150 .5    100 1    230 1  100 .4  445 0.01 ;exponential shark fin
f6 0 1025 5 .01 1025 .01
f7 0 1025 5 .01 200 .5    100 1  300 1   212 .5   200 .3  313 .01 ;throb


;         freq  amp simamp ifn idel   carifn  modfreq modamp modifn pulseNVamp pulseshape minpulse pulsehigh simthr
;         p4      5   6    7    8       9       10     11     12      13        14       15         16       17
i1 0 3000 146.83 .015 1   1    .8     1       97.8     25     1     1           5       .8         1.5      .6
i1 0.153 3000 146.83 .015 1   1    .68     1       97.8     25     1     1           5       .8         1.5      .6


</CsScore>
</CsoundSynthesizer>

<bsbPanel>
 <label>Widgets</label>
 <objectName/>
 <x>100</x>
 <y>100</y>
 <width>320</width>
 <height>240</height>
 <visible>true</visible>
 <uuid/>
 <bgcolor mode="nobackground">
  <r>255</r>
  <g>255</g>
  <b>255</b>
 </bgcolor>
</bsbPanel>
<bsbPresets>
</bsbPresets>
