package play

import (
	"testing"

	"github.com/chrisconley/go-hangman/internal/game"
	"github.com/chrisconley/go-hangman/internal/utils"
	"github.com/chrisconley/go-hangman/test"
)

func TestCountUniqueCharacters(t *testing.T) {
	actual := CountUniqueCharacters("aab")
	test.AssertDeepEqual(actual, 2, t)

	actual = CountUniqueCharacters("cba")
	test.AssertDeepEqual(actual, 3, t)

	actual = CountUniqueCharacters("ccbbaadd")
	test.AssertDeepEqual(actual, 4, t)
}

func TestCountHeldCharacters(t *testing.T) {
	actual := CountHeldCharacters("abc", "abc")
	test.AssertDeepEqual(actual, 3, t)

	actual = CountHeldCharacters("abc", "cba")
	test.AssertDeepEqual(actual, 1, t)

	actual = CountHeldCharacters("abc", "cab")
	test.AssertDeepEqual(actual, 0, t)
}

func TestPlayMastermindSymmetry3(t *testing.T) {
	reactions := []game.Reaction{
		game.NewReaction(game.NewAction("ABCD"), "BW", 0.0),
	}

	lines := utils.ReadLinesFromFilePath("../../fixtures/commercial-mastermind.txt")
	game.BUniverse = game.InitUniverse(lines)
	responseCounts := NewPotentialSymmetries(reactions)
	groupedActions := responseCounts.GroupEquivalentActions()
	expected := [][]string{
		{"AAAA", "BBBB", "CCCC", "DDDD"},
		{"AAAB", "AAAC", "AABA", "AADA", "ACAA", "ADAA", "BBAB", "BBBA", "BBBC", "BBDB", "BCCC", "BDDD", "CACC", "CBBB", "CCCA", "CCCB", "CDCC", "CDDD", "DADD", "DBBB", "DCCC", "DCDD", "DDAD", "DDBD"},
		{"AAAD", "AACA", "ABAA", "ABBB", "ACCC", "ADDD", "BBBD", "BBCB", "CBCC", "CCCD", "DBDD", "DDCD"},
		{"AAAE", "AAAF", "AAEA", "AAFA", "AEAA", "AFAA", "BBBE", "BBBF", "BBEB", "BBFB", "CCCE", "CCCF", "CECC", "CFCC", "DDED", "DDFD", "DEDD", "DFDD", "EBBB", "ECCC", "EDDD", "FBBB", "FCCC", "FDDD"},
		{"AABB", "ACAC", "ADDA", "BBAA", "BCCB", "BDBD", "CACA", "CBBC", "CCDD", "DAAD", "DBDB", "DDCC"},
		{"AABC", "AADB", "ACAB", "ACDA", "ADAC", "ADBA", "BBAC", "BBDA", "BCCA", "BCDD", "BDAD", "BDCC", "CACB", "CADD", "CBBA", "CBDB", "CDBD", "CDCA", "DABD", "DACC", "DBAB", "DBBC", "DCAD", "DCCB"},
		{"AABD", "AACB", "ABAC", "ABBC", "ABDA", "ABDB", "ACAD", "ACCB", "ACDD", "ADBD", "ADCA", "ADCC", "BBAD", "BBCA", "BCCD", "BDCD", "CACD", "CBBD", "CBCA", "CBDD", "DACD", "DBAD", "DBCB", "DBCC"},
		{"AABE", "AABF", "AAEB", "AAFB", "ACAE", "ACAF", "ADEA", "ADFA", "AEAC", "AEDA", "AFAC", "AFDA", "BBAE", "BBAF", "BBEA", "BBFA", "BCCE", "BCCF", "BDED", "BDFD", "CACE", "CACF", "CBBE", "CBBF", "CECA", "CEDD", "CFCA", "CFDD", "DAED", "DAFD", "DBEB", "DBFB", "DEAD", "DECC", "DFAD", "DFCC", "EBBC", "EBDB", "ECCB", "ECDD", "EDBD", "EDCC", "FBBC", "FBDB", "FCCB", "FCDD", "FDBD", "FDCC"},
		{"AACC", "AADD", "ABAB", "ABBA", "ACCA", "ADAD", "BBCC", "BBDD", "CBCB", "CDCD", "DBBD", "DCCD"},
		{"AACD", "ABAD", "ABBD", "ABCA", "ABCB", "ABCC", "ABDD", "ACCD", "ADCD", "BBCD", "CBCD", "DBCD"},
		{"AACE", "AACF", "AAED", "AAFD", "ABAE", "ABAF", "ABBE", "ABBF", "ABEA", "ABEB", "ABFA", "ABFB", "ACCE", "ACCF", "ADED", "ADFD", "AEAD", "AECA", "AECC", "AEDD", "AFAD", "AFCA", "AFCC", "AFDD", "BBCE", "BBCF", "BBED", "BBFD", "CBCE", "CBCF", "CECD", "CFCD", "DBED", "DBFD", "DECD", "DFCD", "EBBD", "EBCB", "EBCC", "EBDD", "ECCD", "EDCD", "FBBD", "FBCB", "FBCC", "FBDD", "FCCD", "FDCD"},
		{"AADC", "ACBA", "ADAB", "BACC", "BADD", "BBDC", "CBAB", "CDAD", "CDCB", "DBBA", "DCBD", "DCCA"},
		{"AADE", "AADF", "AAEC", "AAFC", "ACEA", "ACFA", "ADAE", "ADAF", "AEAB", "AEBA", "AFAB", "AFBA", "BBDE", "BBDF", "BBEC", "BBFC", "BECC", "BEDD", "BFCC", "BFDD", "CBEB", "CBFB", "CDCE", "CDCF", "CDED", "CDFD", "CECB", "CFCB", "DBBE", "DBBF", "DCCE", "DCCF", "DCED", "DCFD", "DEBD", "DFBD", "EACC", "EADD", "EBAB", "EBBA", "ECCA", "EDAD", "FACC", "FADD", "FBAB", "FBBA", "FCCA", "FDAD"},
		{"AAEE", "AAFF", "AEAE", "AEEA", "AFAF", "AFFA", "BBEE", "BBFF", "CECE", "CFCF", "DEED", "DFFD", "EBBE", "EBEB", "ECCE", "EDED", "EECC", "EEDD", "FBBF", "FBFB", "FCCF", "FDFD", "FFCC", "FFDD"},
		{"AAEF", "AAFE", "AEAF", "AEFA", "AFAE", "AFEA", "BBEF", "BBFE", "CECF", "CFCE", "DEFD", "DFED", "EBBF", "EBFB", "ECCF", "EDFD", "EFCC", "EFDD", "FBBE", "FBEB", "FCCE", "FDED", "FECC", "FEDD"},
		{"ABCD"},
		{"ABCE", "ABCF", "ABED", "ABFD", "AECD", "AFCD", "EBCD", "FBCD"},
		{"ABDC", "ACBD", "ADCB", "BACD", "CBAD", "DBCA"},
		{"ABDE", "ABDF", "ABEC", "ABFC", "ACED", "ACFD", "ADCE", "ADCF", "AEBD", "AECB", "AFBD", "AFCB", "BECD", "BFCD", "CBED", "CBFD", "DBCE", "DBCF", "EACD", "EBAD", "EBCA", "FACD", "FBAD", "FBCA"},
		{"ABEE", "ABFF", "AECE", "AEED", "AFCF", "AFFD", "EBCE", "EBED", "EECD", "FBCF", "FBFD", "FFCD"},
		{"ABEF", "ABFE", "AECF", "AEFD", "AFCE", "AFED", "EBCF", "EBFD", "EFCD", "FBCE", "FBED", "FECD"},
		{"ACBB", "ACBC", "ACDC", "ADBB", "ADDB", "ADDC", "BAAD", "BABD", "BACA", "BACB", "BCBD", "BDCB", "CAAD", "CBAA", "CBAC", "CBDC", "CCAD", "CCBD", "DACA", "DBAA", "DBDA", "DBDC", "DDCA", "DDCB"},
		{"ACBE", "ACBF", "ADEB", "ADFB", "AEDC", "AFDC", "BACE", "BACF", "BAED", "BAFD", "CBAE", "CBAF", "CEAD", "CFAD", "DBEA", "DBFA", "DECA", "DFCA", "EBDC", "ECBD", "EDCB", "FBDC", "FCBD", "FDCB"},
		{"ACDB", "ADBC", "BCAD", "BDCA", "CABD", "CBDA", "DACB", "DBAC"},
		{"ACDE", "ACDF", "ACEB", "ACFB", "ADBE", "ADBF", "ADEC", "ADFC", "AEBC", "AEDB", "AFBC", "AFDB", "BCED", "BCFD", "BDCE", "BDCF", "BEAD", "BECA", "BFAD", "BFCA", "CAED", "CAFD", "CBDE", "CBDF", "CBEA", "CBFA", "CEBD", "CFBD", "DACE", "DACF", "DBAE", "DBAF", "DBEC", "DBFC", "DECB", "DFCB", "EABD", "EACB", "EBAC", "EBDA", "ECAD", "EDCA", "FABD", "FACB", "FBAC", "FBDA", "FCAD", "FDCA"},
		{"ACEC", "ACFC", "ADDE", "ADDF", "AEBB", "AFBB", "BEBD", "BECB", "BFBD", "BFCB", "CBEC", "CBFC", "CCED", "CCFD", "DBDE", "DBDF", "DDCE", "DDCF", "EAAD", "EACA", "EBAA", "FAAD", "FACA", "FBAA"},
		{"ACEE", "ACFF", "ADEE", "ADFF", "AEBE", "AEDE", "AEEB", "AEEC", "AFBF", "AFDF", "AFFB", "AFFC", "BECE", "BEED", "BFCF", "BFFD", "CBEE", "CBFF", "CEED", "CFFD", "DBEE", "DBFF", "DECE", "DFCF", "EACE", "EAED", "EBAE", "EBDE", "EBEA", "EBEC", "ECED", "EDCE", "EEAD", "EEBD", "EECA", "EECB", "FACF", "FAFD", "FBAF", "FBDF", "FBFA", "FBFC", "FCFD", "FDCF", "FFAD", "FFBD", "FFCA", "FFCB"},
		{"ACEF", "ACFE", "ADEF", "ADFE", "AEBF", "AEDF", "AEFB", "AEFC", "AFBE", "AFDE", "AFEB", "AFEC", "BECF", "BEFD", "BFCE", "BFED", "CBEF", "CBFE", "CEFD", "CFED", "DBEF", "DBFE", "DECF", "DFCE", "EACF", "EAFD", "EBAF", "EBDF", "EBFA", "EBFC", "ECFD", "EDCF", "EFAD", "EFBD", "EFCA", "EFCB", "FACE", "FAED", "FBAE", "FBDE", "FBEA", "FBEC", "FCED", "FDCE", "FEAD", "FEBD", "FECA", "FECB"},
		{"AEEE", "AFFF", "EBEE", "EECE", "EEED", "FBFF", "FFCF", "FFFD"},
		{"AEEF", "AEFE", "AEFF", "AFEE", "AFEF", "AFFE", "EBEF", "EBFE", "EBFF", "EECF", "EEFD", "EFCE", "EFCF", "EFED", "EFFD", "FBEE", "FBEF", "FBFE", "FECE", "FECF", "FEED", "FEFD", "FFCE", "FFED"},
		{"BAAA", "BABB", "BCBB", "BDBB", "CAAA", "CCAC", "CCBC", "CCDC", "DAAA", "DDDA", "DDDB", "DDDC"},
		{"BAAB", "BABA", "BCBC", "BDDB", "CAAC", "CCAA", "CCBB", "CDDC", "DADA", "DCDC", "DDAA", "DDBB"},
		{"BAAC", "BABC", "BADA", "BADB", "BCBA", "BCDC", "BDAB", "BDDC", "CAAB", "CADC", "CCAB", "CCBA", "CDAA", "CDAC", "CDBB", "CDDB", "DABA", "DADC", "DCAA", "DCBB", "DCBC", "DCDA", "DDAB", "DDBA"},
		{"BAAE", "BAAF", "BABE", "BABF", "BAEA", "BAEB", "BAFA", "BAFB", "BCBE", "BCBF", "BDEB", "BDFB", "CAAE", "CAAF", "CCAE", "CCAF", "CCBE", "CCBF", "CEAA", "CEAC", "CEDC", "CFAA", "CFAC", "CFDC", "DAEA", "DAFA", "DDEA", "DDEB", "DDFA", "DDFB", "DEAA", "DEDA", "DEDC", "DFAA", "DFDA", "DFDC", "ECBB", "ECBC", "ECDC", "EDBB", "EDDB", "EDDC", "FCBB", "FCBC", "FCDC", "FDBB", "FDDB", "FDDC"},
		{"BADC", "CDAB", "DCBA"},
		{"BADE", "BADF", "BAEC", "BAFC", "BEDC", "BFDC", "CDAE", "CDAF", "CDEB", "CDFB", "CEAB", "CFAB", "DCBE", "DCBF", "DCEA", "DCFA", "DEBA", "DFBA", "EADC", "ECBA", "EDAB", "FADC", "FCBA", "FDAB"},
		{"BAEE", "BAFF", "CEAE", "CFAF", "DEEA", "DFFA", "ECBE", "EDEB", "EEDC", "FCBF", "FDFB", "FFDC"},
		{"BAEF", "BAFE", "CEAF", "CFAE", "DEFA", "DFEA", "ECBF", "EDFB", "EFDC", "FCBE", "FDEB", "FEDC"},
		{"BCAA", "BCAB", "BCAC", "BCDB", "BDAA", "BDBA", "BDBC", "BDDA", "CABA", "CABB", "CABC", "CADA", "CCDA", "CCDB", "CDBC", "CDDA", "DAAB", "DAAC", "DABB", "DADB", "DCAC", "DCDB", "DDAC", "DDBC"},
		{"BCAE", "BCAF", "BDEA", "BDFA", "CABE", "CABF", "CEDA", "CFDA", "DAEB", "DAFB", "DEAC", "DFAC", "ECDB", "EDBC", "FCDB", "FDBC"},
		{"BCDA", "BDAC", "CADB", "CDBA", "DABC", "DCAB"},
		{"BCDE", "BCDF", "BCEA", "BCFA", "BDAE", "BDAF", "BDEC", "BDFC", "BEAC", "BEDA", "BFAC", "BFDA", "CADE", "CADF", "CAEB", "CAFB", "CDBE", "CDBF", "CDEA", "CDFA", "CEBA", "CEDB", "CFBA", "CFDB", "DABE", "DABF", "DAEC", "DAFC", "DCAE", "DCAF", "DCEB", "DCFB", "DEAB", "DEBC", "DFAB", "DFBC", "EABC", "EADB", "ECAB", "ECDA", "EDAC", "EDBA", "FABC", "FADB", "FCAB", "FCDA", "FDAC", "FDBA"},
		{"BCEB", "BCFB", "BDBE", "BDBF", "BEAA", "BFAA", "CAEA", "CAFA", "CCDE", "CCDF", "CEBC", "CFBC", "DAAE", "DAAF", "DDEC", "DDFC", "DEDB", "DFDB", "EABB", "ECAC", "EDDA", "FABB", "FCAC", "FDDA"},
		{"BCEC", "BCFC", "BDDE", "BDDF", "BEAB", "BEBA", "BEBC", "BEDB", "BFAB", "BFBA", "BFBC", "BFDB", "CAEC", "CAFC", "CCEA", "CCEB", "CCFA", "CCFB", "CDDE", "CDDF", "CDEC", "CDFC", "CEBB", "CFBB", "DADE", "DADF", "DCDE", "DCDF", "DCEC", "DCFC", "DDAE", "DDAF", "DDBE", "DDBF", "DEBB", "DFBB", "EAAB", "EAAC", "EABA", "EADA", "ECAA", "EDAA", "FAAB", "FAAC", "FABA", "FADA", "FCAA", "FDAA"},
		{"BCEE", "BCFF", "BDEE", "BDFF", "BEAE", "BEEA", "BFAF", "BFFA", "CAEE", "CAFF", "CEBE", "CEDE", "CEEA", "CFBF", "CFDF", "CFFA", "DAEE", "DAFF", "DEAE", "DEEB", "DEEC", "DFAF", "DFFB", "DFFC", "EABE", "EAEB", "ECAE", "ECDE", "ECEB", "EDBE", "EDEA", "EDEC", "EEAC", "EEBC", "EEDA", "EEDB", "FABF", "FAFB", "FCAF", "FCDF", "FCFB", "FDBF", "FDFA", "FDFC", "FFAC", "FFBC", "FFDA", "FFDB"},
		{"BCEF", "BCFE", "BDEF", "BDFE", "BEAF", "BEFA", "BFAE", "BFEA", "CAEF", "CAFE", "CEBF", "CEDF", "CEFA", "CFBE", "CFDE", "CFEA", "DAEF", "DAFE", "DEAF", "DEFB", "DEFC", "DFAE", "DFEB", "DFEC", "EABF", "EAFB", "ECAF", "ECDF", "ECFB", "EDBF", "EDFA", "EDFC", "EFAC", "EFBC", "EFDA", "EFDB", "FABE", "FAEB", "FCAE", "FCDE", "FCEB", "FDBE", "FDEA", "FDEC", "FEAC", "FEBC", "FEDA", "FEDB"},
		{"BEBB", "BFBB", "CCEC", "CCFC", "DDDE", "DDDF", "EAAA", "FAAA"},
		{"BEBE", "BEEB", "BFBF", "BFFB", "CCEE", "CCFF", "CEEC", "CFFC", "DDEE", "DDFF", "DEDE", "DFDF", "EAAE", "EAEA", "ECEC", "EDDE", "EEAA", "EEBB", "FAAF", "FAFA", "FCFC", "FDDF", "FFAA", "FFBB"},
		{"BEBF", "BEFB", "BFBE", "BFEB", "CCEF", "CCFE", "CEFC", "CFEC", "DDEF", "DDFE", "DEDF", "DFDE", "EAAF", "EAFA", "ECFC", "EDDF", "EFAA", "EFBB", "FAAE", "FAEA", "FCEC", "FDDE", "FEAA", "FEBB"},
		{"BEDE", "BEEC", "BFDF", "BFFC", "CDEE", "CDFF", "CEEB", "CFFB", "DCEE", "DCFF", "DEBE", "DFBF", "EADE", "EAEC", "ECEA", "EDAE", "EEAB", "EEBA", "FADF", "FAFC", "FCFA", "FDAF", "FFAB", "FFBA"},
		{"BEDF", "BEFC", "BFDE", "BFEC", "CDEF", "CDFE", "CEFB", "CFEB", "DCEF", "DCFE", "DEBF", "DFBE", "EADF", "EAFC", "ECFA", "EDAF", "EFAB", "EFBA", "FADE", "FAEC", "FCEA", "FDAE", "FEAB", "FEBA"},
		{"BEEE", "BFFF", "CEEE", "CFFF", "DEEE", "DFFF", "EAEE", "ECEE", "EDEE", "EEAE", "EEBE", "EEDE", "EEEA", "EEEB", "EEEC", "FAFF", "FCFF", "FDFF", "FFAF", "FFBF", "FFDF", "FFFA", "FFFB", "FFFC"},
		{"BEEF", "BEFE", "BFEF", "BFFE", "CEEF", "CEFF", "CFEE", "CFFE", "DEFE", "DEFF", "DFEE", "DFEF", "EAEF", "EAFE", "ECEF", "ECFF", "EDFE", "EDFF", "EEAF", "EEBF", "EEFA", "EEFB", "EFAE", "EFBF", "EFDE", "EFDF", "EFEA", "EFEC", "EFFB", "EFFC", "FAEF", "FAFE", "FCEE", "FCFE", "FDEE", "FDEF", "FEAF", "FEBE", "FEDE", "FEDF", "FEEB", "FEEC", "FEFA", "FEFC", "FFAE", "FFBE", "FFEA", "FFEB"},
		{"BEFF", "BFEE", "CEFE", "CFEF", "DEEF", "DFFE", "EAFF", "ECFE", "EDEF", "EEDF", "EEFC", "EFAF", "EFBE", "EFEB", "EFFA", "FAEE", "FCEF", "FDFE", "FEAE", "FEBF", "FEEA", "FEFB", "FFDE", "FFEC"},
		{"EEEE", "FFFF"},
		{"EEEF", "EEFE", "EFEE", "EFFF", "FEEE", "FEFF", "FFEF", "FFFE"},
		{"EEFF", "EFEF", "EFFE", "FEEF", "FEFE", "FFEE"},
	}
	test.AssertDeepEqual(expected, groupedActions, t)
}
