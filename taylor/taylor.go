package taylor

import (
	"math/rand"
	"time"

	"github.com/bclicn/color"
)

var headshots []string

func init() {
	headshots = append(headshots, `
  ..:,,,.:?NN8$O7+~+I7=++=I7I7Z7:I?~.   
  ......,:MDZZ$=~~~:?$O$D7$OZ8ZZ==7I... 
.....,,,,MN8OOO8NNNNDNNN8D8O8D88O$I+++,.
.,,+..:?:?ZOI?=~~:,,,:+77ZZIII?7?7=~~++.
..,,.:=+..Z$$?~,,,,,,,,~+77I=:~??Z7:~:+?
..,,,:...$IO+=+~:,,,,,,,:==+=::~+77+=~7I
.,,~...:7==$$+~::~++,,,,,~77I+==~+?:=++.
,,,...,D$=?ZOI$NMDI?=,,,:~ONMN7OO7=+++=.
,,,,..~$I777~:=??+=~~:,::==Z~$II7??+I~,.
...+8O7+:+?=::,,,,,,~:,::::::~:~7$+~~~..
.7I+DZ+:,=?=~~:,,,,:~:,:~,,,:::~+?=,:::.
.+??OO7$,=7?=~::,,,~~,,:~:,,::~=I=~:=::.
.,:=?=~~?+II+=~:,,,=$~:~::,:::=+7=+=::~.
.,,~==+:=?,~I=~:::,::::::,,::~7:,:~:,~=~
,,,::+~:=:.:=?~:::::==~=:,:::=~,.,~,,?=,
,,,,:~+:::,,=Z+~::~?=++~+~::=+$I,,+:.:::
,,,,::~~~,,.,,:7~:::=++:::~+D+I:,,,:,:~=
,:~I::=,,~:,,,.,,=~:,,,,::+::,,,,,,~,,.:
+:7Z~:,~~:~~::,,.,==:~~~=:,,,,,:~+:..,..
~7?ZO=+=~~~~~~::,,,,IO+:,,,,:~=++++::.,,
~=:7+8ODI?=~~:::,,,,,8~,,,::==8?I?==~~=:
~~$Z~IZODOO$I=~~:,,,.::,,,:~$DZ$7?=~~?~,
~~I:=?~?I888OO+=~,.,,,:,,,,,NN=$=?=+?+:,
~~=~+?::?:~IOOO+~:,,,,+:,,,:?$=??==??~:,
~~~I=:II7~:::7Z=~:,,,,~:,,,:~7:?,::++~=:
`)

	headshots = append(headshots, `
???++=++O88888NDNN8D=ODND8888ZZOO8OO$:::
??++++===Z8DNNND=.=......8888888888O$~::
??+++++...ZDNNO,7Z$OZIZI..+ZO888888OZI~:
??++=+I~:ZDNNN~$??++8N7DO..ZZOZ8888OO$=:
??+++++?ONNNN8$7??=I==?7?~=:Z$888O88OZ::
??+++==8:DNNNDOZ+=7?+?7ZII~.$888OZOOOO:.
?++++=?DNDD$$?+Z8+Z$$?IO87:?$8OZZZZ$$=.:
??++OO88ZZD8D8I$II++7ZO8O~=,:IZZZZ$7:...
??+?O88888OOD$Z+ZOO$ZO8$7?,:,.OZZ$78~...
+?+OOOOO7?8NNZ=+7DD?+7+7I7+~7,OZ8Z8OZZO8
~=IOOOOOOOZII8I+?+I++++:I$=:.,$.ODD8ODDN
~=Z,:,ZZO8DD$NNNDN:=:~,::=~=:,:~..ND8IDD
.~?::OOOO88DDN$ :,,=~7:~~==~=:.:..:.NNDZ
8~+:.OOO8Z8$Z,,Z,,~:,+:+=~~?=I=7~,:?.NNN
~?~=OOOODD8Z~:.:====::I~~?+~I++=,,,~,.NN
?+~NNZZ878I,~==:=~~=~~,:=?+?7I??+~==::,N
+,~.88DD++..,:~~:+=~=Z==+=I?$?7NNN:+=,,+
+~,:$8Z7+7:=~==+?+~+O~:?=+?I=7IN8ONI~~:~
$~::.ZZ7II:~+=+=+~~~,=,+~+++7?I8NNNN?O,=
I+:,:.Z7~,=~+~?=++++++=+?II7ZIDNNNND7+~?
II=8=:.,:::~~N??++++++++?I777+NNN8???~7+
II?~:::?,:~7ONI?+++==+??I7777NNNN888~::=
?I??+=:::7,7O8I?+++++?I?III$MD88D88::,=7
.~O8+7+=~~DDD7I?+?+++??++??7DD8888~:=+IO
O8DDN+=ID$D8===$+~:?+?7Z+?++78+~Z+Z8?Z88
`)
	headshots = append(headshots, `
I?+7$$$$$$$$$$$$$$$$$$$ZZ$ZZOOOO888888DD
7II7$$$$$$$$$$$I~$O88$7ZZZZZOOOOO8888DDD
$77$Z$$$$$$$$$?$O$$$$DN8I$ZZOOO888888DDD
$7$ZZ$$$$$$$+$ZZZZOO$$O8O$OOOOOO88888DDD
$7$ZZZZ$$$$+$$$$$$$$$$ZO8ZIOOOO888888DDD
7I$ZZZZZ$$,$$$$$$$$$$ZZZOOZ=OO8888888DDD
IIZZZZZ$Z,$7$Z$$$ZZZZZZZZ8O$$O8888888DDD
IIZZZZZZ?$77$$$Z.O+ZZZZZOO+Z7O8888888DDD
7IZZZZZ,$7$Z$$.~=:?==I$ZZOZO$I8O888888DD
77ZZZZZO$$Z7.:=$$77$$:$ZZOZ7O77O888888DD
7$OOZZZOZ7Z=,~O$$ZZZ7?.$ZOZ,OZZZ88888D8D
7$OOOZZ8Z$$=:7D$7I$IIO:,ZO,ZZZ8O88888DDD
7$OOOOZ8O$$,IZO$$$D8::~$Z8ZZZZ8O8888DDDD
7ZOOOOZ$O$$=:O8OZ$$7OZ.:$$ZZZ88O8888DDDD
7ZOOOOOOOZ$:8$DZO88DZ+:Z$ZZZO88O8888DDDD
$OOOOOOO8OZ7:OZO7$78~?I7ZOZOO88O888DDDDD
$OOOOOOOZO$7Z8$$7I77NN+OZZOOO88O888DDDDD
`)

	headshots = append(headshots, `
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMM~~~~~~~~~~~~~~~MMMMMMMMMMMM
MMMMMMMMM8~~~~~~~~~~~~~~~~~~~~MMMMMMMMMM
MMMMMMN:N~~~~~~~~~~~~~~~~~~~~~NNNMMMMMMM
MMMMMM:~=~~~~~~~~~~~~~~~~~~~:=~::MMMMMMM
MMMMM:N~~~~=~N~~~~~~~~~~~~~~~~~=:~MMMMMM
MMM++++7????77???~~~:~??????I$$IIIIIMMMM
MMM++I7??7?7$$$$+?????77??7?7$$$$IIIMMMM
MMM:+7????$$$$$$$?::??7????777$$$$I~MMMM
MMM~+7????$$$$$$$?:::?7????777$$$$I~:MMM
MMM~+??7?$$$$$$$?:NZN???7?7777$$$?~~~MMM
MMM~~+??77$$$$$$?INN==???77777$$$I~~~MMM
MMM~~~+?777$$$?:NIIDI~M??7777$$$I~~~~MMM
MMM~~~::+?++::::::+DN:::::=III~::~~~~MMM
MMM~~~:::::::::::::::::::::::::::~~~~MMM
MMM~~~:::::,???????++++++::::::::~~~~MMM
MMM~~~I????????????+++++++++++++:~~~~MMM
MMMM~I~????????????+++++++++++++~~~~MMMM
MMMM~~??????????????++++++++++++~+~~MMMM
MMMMM~??????????????+++++++++++~++~MMMMM
MMMMM:??????????????++++++++++++++~MMMMM
MMMMM???????????????++++++++++++++MMMMMM
MMMMM??????????????I++++++++++++++MMMMMM
`)

	headshots = append(headshots, `
DDDDDDDDDDDD+D,I=8DDD888DD88DD8DDDDDDDD8
DDDDDDDDDD=8$OZ$O7?+.8888888888DDDDD8D88
DDDDDDDDDIZ8?DDD8$??I?=+8888888888888D88
DDDDDDDN+$7?=+8IZZZZ$I+=~888888D8888DDDD
DDDDDD=?ZZZDDD8D8888888+==.88888DDDDDDDD
DDDDDD=?IDDDDDDD88888888?===+8888DDDDDDD
DDDDD:+I8DDDDDDDDDD8888887===.8888DDDDDD
DDDDD++IDDDDDDDDDD88888888O+===:888DDDDD
DDDD++?=DDDDDDD8D88888888888?===,88DDDDD
DDD8=+I$DDDDDD88:.....O888888?+=$D888DD8
DDD+=?I8DDDDD8.,.,+$:....O888++??8D8DDD8
DD8==?+DDDD8?.=7??77$O~,.,Z88+++?:DDDDD8
DD+=+I:DDDDO:=$+===+I?=?=,:88+++I?DDDDD8
DDZI??.DDDDI.?O+OD=+?+II$I~+8++??7DDDDD8
DDDII??:DDZ.,$8+==?7+?I$$=,:,:??I7DDDDD8
DDDDII???8DI.$8O?I+O+??IOZ..?I??I?DDDDDD
DDDDD7I???DI=?ZDN??I=+IZZDDN:+?IIZDDDDDD
DDDDDD7I???$,=IODD$$$O7Z8+I8.??II8DDDDDD
DDDDDDD7III?+,+$ZD$I77?I7I$7I???IDDDDDDD
DDDDDDDD7II??=?:=77?+?+?$OOOI?+?IDDDDDDD
DDDDDDDDD77?++=+O.==+==++DZDO7++7DDDDDDD
DDDDDDDDDD77?+=+8+=======N:N87++8DDDDDDD
DDDDDDDDDDDNI+?ZZO++++===?NOD77?DDDDDDDD
DDDDDDDDDDDDD?+ZOO8D====+ZD?DO$?DDDDDDDD
DDDDDDDDDDDDD8+I8OD88IO7IODZDI8$DDDDDDDD
`)
}

// HeadShot returns ascii-art for a face shot
func HeadShot() string {
	rand.Seed(time.Now().Unix())
	return color.Purple(headshots[rand.Intn(len(headshots))])
}
