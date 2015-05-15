package main

import (
	"flag"
	"fmt"
	"github.com/huichen/sego"
	"github.com/jbrukh/bayesian"
	"os"
)

var brainFile *string = flag.String("brain", "0.bra", "Brain File Path")
var input *string = flag.String("input", "", "Your input")

const (
	Good bayesian.Class = "Good"
	Bad  bayesian.Class = "Bad"
)

func createClassifier() *bayesian.Classifier {
	flag.Parse()
	var classifier *bayesian.Classifier
	if _, err := os.Stat(*brainFile); err == nil {
		fmt.Printf("Load Brain from %v\n", *brainFile)
		classifier, _ = bayesian.NewClassifierFromFile(*brainFile)
	} else {
		fmt.Printf("No Brain file detected, create one\n")
		classifier = bayesian.NewClassifier(Good, Bad)
		goodStuff := []string{"高", "有钱", "帅"}
		badStuff := []string{"穷", "脏", "丑"}
		classifier.Learn(goodStuff, Good)
		classifier.Learn(badStuff, Bad)
		classifier.WriteToFile("./1.bra")
	}
	return classifier
}

func createSegmenter() *sego.Segmenter {
	segmenter := new(sego.Segmenter)
	segmenter.LoadDictionary("./dictionary.txt")
	// 分词
	// text := []byte("中华人民共和国中央人民政府")
	// segments := segmenter.Segment(text)
	// fmt.Println(sego.SegmentsToString(segments, false))

	// text = []byte("本章节描述了本文档的发布状态")
	// segments = segmenter.Segment(text)

	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	// fmt.Println(sego.SegmentsToString(segments, false))
	return segmenter
}

func main() {
	segmenter := createSegmenter()
	classifier := createClassifier()

	var probs []float64
	var likely int
	// var err error
	var words []string

	if *input != "" {
		segments := segmenter.Segment([]byte(*input))
		words = sego.SegmentsToSlice(segments, false)
		// fmt.Printf("%v\n", words)
	}

	// probs, likely, _, err = classifier.SafeProbScores([]string{"高个", "妹子"})
	// if err == nil {
	// 	fmt.Printf("Good=%v, Bad=%v, likely=%v\n", probs[0], probs[1], likely)
	// }

	probs, likely, _ = classifier.ProbScores(words)
	// probs, likely, _, err = classifier.SafeProbScores([]string{"一个", "有钱", "的", "丑", "妹子"})
	fmt.Printf("Good=%v, Bad=%v, likely=%v\n", probs[0], probs[1], likely)
	if probs[0] < 0.1 {
		fmt.Printf("不建议选择 %v。\n", *input)
	} else if probs[0] > 0.8 {
		fmt.Printf("建议选择 %v。\n", *input)
	} else {
		fmt.Printf("我也搞不清楚是否该选择 %v，建议你自己看着办。\n", *input)
	}
	return
}
