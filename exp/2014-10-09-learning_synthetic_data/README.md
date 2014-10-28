Synthesize the corpus:

  $GOPATH/bin/synthesize -model=ground_truth_model.json -corpus=synthetic.json -instances=10000 -cardi=1 -length=10

Train the model:

  $GOPATH/bin/trainer -corpus=synthetic.json -model=learned_model.json -iter=200 -logl=logl.txt -states=10