<CsoundSynthesizer>
<CsOptions>

</CsOptions>
<CsInstruments>

sr = 44100
ksmps = 100 ;samps / k-period
nchnls = 2
;0dbfs = 1.0

;todo: reverb and/or delay
gihandle OSCinit 7770
;ga1 init 0

instr 1
    kCount init 0
	kf1 init 0
	kf2 init 0
	kf3 init 0
	kf4 init 0
	kf5 init 0
	kf6 init 0
	kf7 init 0
	kf8 init 0
	;kf9 init 0

	ibasefreq = p4
	iscale = p5
	ifn = p6
	idelbase = p7
nxtmsg:

 	kk OSClisten gihandle, "/eeg", "iiiiiiiii", kCount, kf1, kf2, kf3, kf4, kf5, kf6, kf7, kf8

	if (kk == 0) goto no_new_data
	printks "%i %i %i %i %i %i %i %i %i\\n", 0, kCount, kf1, kf2, kf3, kf4, kf5, kf6, kf7, kf8
no_new_data:

 	;USING DATA FOR AMPLITUDE
 	asig1 oscili kf1*iscale, ibasefreq, ifn
 	asig2 oscili kf2*iscale, ibasefreq * 2, ifn
 	asig3 oscili kf3*iscale, ibasefreq * 3, ifn
 	asig4 oscili kf4*iscale, ibasefreq * 4, ifn
 	asig5 oscili kf5*iscale, ibasefreq * 5, ifn
 	asig6 oscili kf6*iscale, ibasefreq * 6, ifn
 	asig7 oscili kf7*iscale, ibasefreq * 7, ifn
 	asig8 oscili kf8*iscale, ibasefreq * 8, ifn
; 	asig9 oscili kf9*iscale, ibasefreq* 9, ifn


 	aoutsig = asig1 + asig2 + asig3 + asig4 + asig5 + asig6 + asig7 + asig8 ; + asig9
 	acomp oscili 1000, 1000, 1
 	abal balance aoutsig, acomp


		afiltsig butlp abal,  1500


kenver oscili  1, 2, 3
aenvfiltsig =  kenver * afiltsig
;ares delay asig, idlt [, iskip]
adelay1 delay aenvfiltsig, idelbase
adelay2 delay aenvfiltsig, idelbase * .2
adelay3 delay aenvfiltsig, idelbase * .3
adelay4 delay aenvfiltsig, idelbase * .4
adelay5 delay aenvfiltsig, idelbase * .5
adelay6 delay aenvfiltsig, idelbase * .6
adelay7 delay aenvfiltsig, idelbase * .7
adelay8 delay aenvfiltsig, idelbase * .8
adelay9 delay aenvfiltsig, idelbase * .9
adelay10 delay aenvfiltsig, idelbase * 1



adelays = adelay1 + adelay2 + adelay3 + adelay4 + adelay5 + adelay6 + adelay7 + adelay8 + adelay9 + adelay10
 abal2 balance adelays, acomp

apanL, apanR pan2 abal2, p8
;apanR pan2 abal2, .7

 	;	ga1 += afiltsig.

 		outs apanL, apanR

	endin

instr 2
ienvpersec = p4
asig oscili 10000, 440, 1
;aenver oscil1 idel, kamp, idur [, ifn]
kenver oscili  1, ienvpersec,  3

outs asig * kenver, asig * kenver
;outs asig , asig

endin


;  instr 3; reverb of gaBus and output
;
;;  ares reverb asig, krvt [, iskip]
;
; igscale = p4
; irevtime = p5
;ares reverb ga1*igscale, irevtime
;
;;aoutL, aoutR reverbsc ainL, ainR, kfblvl, kfco[, israte[, ipitchm[, iskip]]]
;          outs      ares, ares
;  endin
;
;  instr 100; clear global audios at the end
;          clear     ga1, ga1
;  endin


</CsInstruments>
<CsScore>
;f1 0 16384 10 1 0 1 .8 .6 .3 .1
;f1 0 16384 10 1 0 .6 0 .1
f1 0 16384 10 1
f2 0 16384 10 1 0 1 0 .5 0 .2 0 .01
f4 0 16384 10 1 0 1
f3 0 16384 9 0.5 1 0

;              freq scale ifn idelbase   pan
;these two
;i1 0    3000  46     1    1  .0411229     .11
;i1 .2   3000   138  .7   1   .1013524     .85


;these two
i1 .8  3000     36         .7   1  .01   .85
i1 .39 3000     239.4432   .5   1   .002   .15



;i1 0  3000    415  .8   1  .1         .15
;i1 .4  3000     36  .7   1  .121 .85



;i1 .39 2000     184   .5   1   .0254111       .15


;i1 0 3000   55 15 2
;i1 0 3000   82.4 15 2 .7
;i1 .5 40000 123.4 15 1

;          envpersec
;i2 0    1  3

;i s dur scale revtime
;i3 0 20 .1     10
;i100 0 30

;i1 0 3000 277 15 2
;i1 0 3000 415 15 2

;i99 0 4000 .5

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
