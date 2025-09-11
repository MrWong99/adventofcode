package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/MrWong99/adventofcode/pluginmanager/shared"
	"github.com/hashicorp/go-plugin"
)

type CalcService struct{}

func (CalcService) Calculate(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) != 2 {
		return "", fmt.Errorf("expected 2 lines of input, but got %d", len(lines))
	}

	list1, err := parseLine(lines[0])
	if err != nil {
		return "", fmt.Errorf("failed to parse first list: %w", err)
	}

	list2, err := parseLine(lines[1])
	if err != nil {
		return "", fmt.Errorf("failed to parse second list: %w", err)
	}

	if len(list1) != len(list2) {
		return "", fmt.Errorf("the two lists must have the same number of elements, but got %d and %d", len(list1), len(list2))
	}

	sort.Ints(list1)
	sort.Ints(list2)

	var totalDistance int
	for i := 0; i < len(list1); i++ {
		distance := list1[i] - list2[i]
		if distance < 0 {
			distance = -distance
		}
		totalDistance += distance
	}

	return strconv.Itoa(totalDistance), nil
}

func parseLine(line string) ([]int, error) {
	numStrs := strings.Fields(line)
	if len(numStrs) == 0 {
		return []int{}, nil
	}
	nums := make([]int, len(numStrs))
	for i, s := range numStrs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("could not parse '%s' as an integer: %w", s, err)
		}
		nums[i] = n
	}
	return nums, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"day1": &shared.CalcGRPCPlugin{Impl: &CalcService{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
