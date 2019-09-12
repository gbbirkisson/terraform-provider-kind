PLUGIN_DIR:=${HOME}/.terraform.d/plugins
PLUGIN_NAME:=terraform-provider-kind_v0.1.0_x5

install:
	mkdir -p ${PLUGIN_DIR}
	go build -o ${PLUGIN_DIR}/${PLUGIN_NAME} main.go