.:
	go build

clean:
	go clean

install_service:
	sudo cp climateserv.service /lib/systemd/system/ \
		&& sudo chmod 755 /lib/systemd/system/climateserv.service \
		&& sudo systemctl enable climateserv.service \
		&& sudo systemctl start climateserv

disable_service:
	sudo systemctl disable climateserv.service

.PHONY: run
run:
	./climateserv -i 10 -d /dev/ttyUSB0 -p 2510 -i 10
