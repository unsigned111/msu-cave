<CsoundSynthesizer>
<CsOptions>
</CsOptions>
<CsInstruments>

sr = 44100
ksmps = 128
nchnls = 2
0dbfs = 1

giosc1 OSCinit 7770


instr 1
	k1 init 0
	k2 init 0
	k3 init 0
	k4 init 0
	k5 init 0
	k6 init 0
	k7 init 0
	k8 init 0
	kans OSClisten giosc1, "/eeg", "iiiiiiii", k1, k2, k3, k4, k5, k6, k7, k8
	printk 0, k1

	aosc oscil 100, 440+50*(1578604-k2)/1000000, 1
	outs aosc, aosc
endin


</CsInstruments>
<CsScore>
f 1 0 1024 10 1
i 1 0 300
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
