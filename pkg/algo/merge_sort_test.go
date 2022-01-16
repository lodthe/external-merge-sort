package algo

import (
	"bytes"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/lodthe/external-merge-sort/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestMergeSort(t *testing.T) {
	samples := []string{
`
F68C3A0DA6
DCD4BFDDA2
57E28EAD2F
745EB200AA
757AAA4A36
CE7B77E43F
C1CD093AA9
8580385A83`,
`
0833723679
7778318607
0265586573
1713352528
3073801602
5448392084
2488570953
7396157502`,
``,
	}

	for _, sample := range samples {
		checkSample(t, sample)
	}
}

func checkSample(t *testing.T, sample string) {
	file, err := os.CreateTemp("", "test_merge_sort_")
	assert.Nil(t, err, "create input file")
	defer func() {
		_ = os.Remove(file.Name())
	}()

	err = ioutil.WriteFile(file.Name(), []byte(sample), 0644)
	_ = file.Close()
	assert.Nil(t, err, "write input")

	msort := NewExternalMergeSort(&config.Config{
		BlockSize:   2,
		MemoryLimit: 7,
		Delimiter:   byte('\n'),
		Less:        func(a, b []byte) bool {
			return bytes.Compare(a, b) < 0
		},
	})

	const outputFilename = "test_merge_sort_output.txt"
	defer func() {
		_ = os.Remove(outputFilename)
	}()

	err = msort.Sort(file.Name(), outputFilename, ".")
	assert.Nil(t, err, "run sort")

	output, err := ioutil.ReadFile(outputFilename)
	assert.Nil(t, err, "read output")

	sorted := strings.Split(sample, "\n")
	sort.Strings(sorted)

	expected := strings.Join(sorted, "\n")
	if expected != "" {
		expected += "\n"
	}

	assert.Equal(t, expected, string(output), "valid output")
}
