
$GOPATH/bin/synthesize -model=ground_truth_model.json -corpus=synthetic.json -instances=10000 -length=10 -cardi=10 

$GOPATH/bin/trainer -corpus=synthetic.json -model=learned_model.json -iter=200 -logl=logl.txt -states=10