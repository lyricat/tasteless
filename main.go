package main

import (
	"flag"
	"fmt"
	"github.com/huichen/sego"
	"github.com/jbrukh/bayesian"
	"io/ioutil"
	"os"
)

var brainFile *string = flag.String("brain", "0.bra", "Brain File Path")
var input *string = flag.String("input", "", "Your input file")
var learn *bool = flag.Bool("learn", false, "learn")

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
		fmt.Printf("word cound: %v\n", classifier.WordCount())
	} else {
		fmt.Printf("No Brain file detected, create one\n")
		classifier = bayesian.NewClassifier(Good, Bad)
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

func learnFromFile(fileName string) []byte {
	data, err := ioutil.ReadFile(fileName)
	if err == nil {
		return data
	}
	return nil
}

func learnFromSamples(classifier *bayesian.Classifier, segmenter *sego.Segmenter) {
	var content []byte
	var segments []sego.Segment
	var words []string
	var files []os.FileInfo

	files, _ = ioutil.ReadDir("./samples/BAD")
	words = make([]string, 0)
	for _, f := range files {
		content = learnFromFile("samples/BAD/" + f.Name())
		segments = segmenter.Segment(content)
		words = sego.SegmentsToSlice(segments, false)
		fmt.Printf("Learn Bad `%v` (%v words)\n", f.Name(), len(words))
		classifier.Learn(words, Bad)
	}

	files, _ = ioutil.ReadDir("./samples/GOOD")
	words = make([]string, 0)
	for _, f := range files {
		content = learnFromFile("samples/GOOD/" + f.Name())
		segments = segmenter.Segment(content)
		words = sego.SegmentsToSlice(segments, false)
		fmt.Printf("Learn Good `%v` (%v words)\n", f.Name(), len(words))
		classifier.Learn(words, Good)
	}
	return
}

func main() {
	segmenter := createSegmenter()
	classifier := createClassifier()
	if *learn {
		fmt.Println("Learn from samples")
		learnFromSamples(classifier, segmenter)
		classifier.WriteToFile("./0.bra")
		return
	}

	var probs []float64
	var likely int
	// var err error
	var words []string

	if *input != "" {
		data, err := ioutil.ReadFile(*input)
		if err != nil {
			fmt.Printf("Cant read %v, err=%v\n", *input, err)
		} else {
			segments := segmenter.Segment(data)
			words = sego.SegmentsToSlice(segments, false)
			fmt.Printf("Read input, %v words, %v segments\n", len(words), len(segments))
		}
	}

	// probs, likely, _, err = classifier.SafeProbScores([]string{"高个", "妹子"})
	// if err == nil {
	// 	fmt.Printf("Good=%v, Bad=%v, likely=%v\n", probs[0], probs[1], likely)
	// }

	probs, likely, _ = classifier.ProbScores(words)
	// probs, likely, _, err = classifier.SafeProbScores([]string{"一个", "有钱", "的", "丑", "妹子"})
	fmt.Printf("Good=%v, Bad=%v, probs=%v, likely=%v\n", probs[0], probs[1], probs, likely)
	if probs[0] < 0.1 {
		fmt.Printf("%v 很烂\n", *input)
	} else if probs[0] > 0.8 {
		fmt.Printf("%v 不错\n", *input)
	} else {
		fmt.Printf("我也搞不清楚 %v 怎么样，建议你自己看着办。\n", *input)
	}
	return
}
