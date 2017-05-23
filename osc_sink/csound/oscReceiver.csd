<CsoundSynthesizer>
<CsOptions>

; to run, type
; csound -odac +-rtaudio=alsa -B2048 -b2048 oscSendtest2.csd

</CsOptions>
<CsInstruments>

sr = 48000
ksmps = 128
nchnls = 2
0dbfs = 1.0

gihandle OSCinit 7770

instr 1
    ki1 init 1
    ki2 init 2
    ki3 init 3
    ki4 init 4
    ki5 init 5
    ki6 init 6
    ki7 init 7
    ki8 init 8
    ki9 init 9

next_message:
    kans OSClisten gihandle, "/eeg", "iiiiiiiii", ki1, ki2, ki3, ki4, ki5, ki6, ki7, ki8, ki9

    if (kans == 0) goto no_new_data
    printks "%i %i %i %i %i %i %i %i %i\\n", 0, ki1, ki2, ki3, ki4, ki5, ki6, ki7, ki8, ki9

no_new_data:

endin

</CsInstruments>

<CsScore>

f 1 0 1024 10 1
i 1 0 4000

</CsScore>

</CsoundSynthesizer>


