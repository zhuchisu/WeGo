// Copyright © 2017 Makoto Ito
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package builder

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/ynqa/wego/config"
	"github.com/ynqa/wego/model"
	"github.com/ynqa/wego/model/word2vec"
	"github.com/ynqa/wego/validate"
)

// Word2vecBuilder manages the members to build Model interface.
type Word2vecBuilder struct {
	// input file path.
	inputFile string

	// common configs.
	dimension  int
	iteration  int
	minCount   int
	threadSize int
	window     int
	initlr     float64
	toLower    bool
	verbose    bool

	// word2vec configs.
	model              string
	optimizer          string
	batchSize          int
	maxDepth           int
	negativeSampleSize int
	subsampleThreshold float64
	theta              float64
}

// NewWord2vecBuilder creates *Word2vecBuilder.
func NewWord2vecBuilder() *Word2vecBuilder {
	return &Word2vecBuilder{
		inputFile: config.DefaultInputFile,

		dimension:  config.DefaultDimension,
		iteration:  config.DefaultIteration,
		minCount:   config.DefaultMinCount,
		threadSize: config.DefaultThreadSize,
		window:     config.DefaultWindow,
		initlr:     config.DefaultInitlr,
		toLower:    config.DefaultToLower,
		verbose:    config.DefaultVerbose,

		model:              config.DefaultModel,
		optimizer:          config.DefaultOptimizer,
		batchSize:          config.DefaultBatchSize,
		maxDepth:           config.DefaultMaxDepth,
		negativeSampleSize: config.DefaultNegativeSampleSize,
		subsampleThreshold: config.DefaultSubsampleThreshold,
		theta:              config.DefaultTheta,
	}
}

// NewWord2vecBuilderFromViper creates *Word2vecBuilder from viper.
func NewWord2vecBuilderFromViper() *Word2vecBuilder {
	return &Word2vecBuilder{
		inputFile: viper.GetString(config.InputFile.String()),

		dimension:  viper.GetInt(config.Dimension.String()),
		iteration:  viper.GetInt(config.Iteration.String()),
		minCount:   viper.GetInt(config.MinCount.String()),
		threadSize: viper.GetInt(config.ThreadSize.String()),
		window:     viper.GetInt(config.Window.String()),
		initlr:     viper.GetFloat64(config.Initlr.String()),
		toLower:    viper.GetBool(config.ToLower.String()),
		verbose:    viper.GetBool(config.Verbose.String()),

		model:              viper.GetString(config.Model.String()),
		optimizer:          viper.GetString(config.Optimizer.String()),
		batchSize:          viper.GetInt(config.BatchSize.String()),
		maxDepth:           viper.GetInt(config.MaxDepth.String()),
		negativeSampleSize: viper.GetInt(config.NegativeSampleSize.String()),
		subsampleThreshold: viper.GetFloat64(config.SubsampleThreshold.String()),
		theta:              viper.GetFloat64(config.Theta.String()),
	}
}

// InputFile sets input file string.
func (wb *Word2vecBuilder) InputFile(inputFile string) *Word2vecBuilder {
	wb.inputFile = inputFile
	return wb
}

// Dimension sets dimension of word vector.
func (wb *Word2vecBuilder) Dimension(dimension int) *Word2vecBuilder {
	wb.dimension = dimension
	return wb
}

// Iteration sets number of iteration.
func (wb *Word2vecBuilder) Iteration(iter int) *Word2vecBuilder {
	wb.iteration = iter
	return wb
}

// MinCount sets min count.
func (wb *Word2vecBuilder) MinCount(minCount int) *Word2vecBuilder {
	wb.minCount = minCount
	return wb
}

// ThreadSize sets number of goroutine.
func (wb *Word2vecBuilder) ThreadSize(threadSize int) *Word2vecBuilder {
	wb.threadSize = threadSize
	return wb
}

// Window sets context window size.
func (wb *Word2vecBuilder) Window(window int) *Word2vecBuilder {
	wb.window = window
	return wb
}

// Initlr sets initial learning rate.
func (wb *Word2vecBuilder) Initlr(initlr float64) *Word2vecBuilder {
	wb.initlr = initlr
	return wb
}

// ToLower is whether converts the words in corpus to lowercase or not.
func (wb *Word2vecBuilder) ToLower() *Word2vecBuilder {
	wb.toLower = true
	return wb
}

// Verbose sets verbose mode.
func (wb *Word2vecBuilder) Verbose() *Word2vecBuilder {
	wb.verbose = true
	return wb
}

// Model sets model of Word2vec. One of: cbow|skip-gram
func (wb *Word2vecBuilder) Model(model string) *Word2vecBuilder {
	wb.model = model
	return wb
}

// Optimizer sets optimizer of Word2vec. One of: hs|ns
func (wb *Word2vecBuilder) Optimizer(optimizer string) *Word2vecBuilder {
	wb.optimizer = optimizer
	return wb
}

// BatchSize sets batch size to update learning rate.
func (wb *Word2vecBuilder) BatchSize(batchSize int) *Word2vecBuilder {
	wb.batchSize = batchSize
	return wb
}

// MaxDepth sets number of times to track huffman tree.
func (wb *Word2vecBuilder) MaxDepth(maxDepth int) *Word2vecBuilder {
	wb.maxDepth = maxDepth
	return wb
}

// NegativeSampleSize sets number of samples as negative.
func (wb *Word2vecBuilder) NegativeSampleSize(size int) *Word2vecBuilder {
	wb.negativeSampleSize = size
	return wb
}

// SubSampleThreshold sets threshold for subsampling.
func (wb *Word2vecBuilder) SubSampleThreshold(threshold float64) *Word2vecBuilder {
	wb.subsampleThreshold = threshold
	return wb
}

// Theta sets lower limit of learning rate (lr >= initlr * theta).
func (wb *Word2vecBuilder) Theta(theta float64) *Word2vecBuilder {
	wb.theta = theta
	return wb
}

// Build creates model.Model interface.
func (wb *Word2vecBuilder) Build() (model.Model, error) {
	if !validate.FileExists(wb.inputFile) {
		return nil, errors.Errorf("Not such a file %s", wb.inputFile)
	}

	input, err := os.Open(wb.inputFile)
	if err != nil {
		return nil, err
	}

	cnf := model.NewConfig(wb.dimension, wb.iteration, wb.minCount, wb.threadSize, wb.window,
		wb.initlr, wb.toLower, wb.verbose)

	var opt word2vec.Optimizer
	switch wb.optimizer {
	case "hs":
		opt = word2vec.NewHierarchicalSoftmax(wb.maxDepth)
	case "ns":
		opt = word2vec.NewNegativeSampling(wb.negativeSampleSize)
	default:
		return nil, errors.Errorf("Invalid optimizer: %s not in hs|ns", wb.optimizer)
	}

	var mod word2vec.Model
	switch wb.model {
	case "cbow":
		mod = word2vec.NewCbow(wb.dimension, wb.window, wb.threadSize)
	case "skip-gram":
		mod = word2vec.NewSkipGram(wb.dimension, wb.window, wb.threadSize)
	default:
		return nil, errors.Errorf("Invalid model: %s not in cbow|skip-gram", wb.model)
	}

	return word2vec.NewWord2vec(input, cnf, mod, opt,
		wb.batchSize, wb.subsampleThreshold, wb.theta)
}
