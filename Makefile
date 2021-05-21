LD_FLAGS="-w -s"

gitea_server_dist=build/server/user_control_by_gitea
file_server_dist=build/server/user_control_by_config_file
cmd_client=build/client/cmd_edge_client

all: base gitea_server config_file_server cmd_client

win: all win_gui_client

gitea_server:
	mkdir $(gitea_server_dist)
	cd server/user_control_by_gitea && go build -ldflags $(LD_FLAGS) -o ../../$(gitea_server_dist)/
	cp server/user_control_by_gitea/config.example.toml $(gitea_server_dist)/config.toml

config_file_server:
	mkdir $(file_server_dist)
	cd server/user_control_by_config_file && go build -ldflags $(LD_FLAGS) -o ../../$(file_server_dist)/
	cp server/user_control_by_config_file/users.toml $(file_server_dist)/

cmd_client:
	mkdir $(cmd_client)
	cd client/cmd_edge_client && go build -ldflags $(LD_FLAGS) -o ../../$(cmd_client)/
	cp client/cmd_edge_client/config.toml $(cmd_client)
	cp client/cmd_edge_client/setup.bat $(cmd_client)
	cp client/cmd_edge_client/setup.sh $(cmd_client)

win_gui_client:


base:
	rm build -rf
	mkdir build
	mkdir build/server
	mkdir build/client

