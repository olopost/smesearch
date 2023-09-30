SEARCHPATH=$(HOME)/bin

install: smesearch
	cp smesearch ~/bin
	cp utils/fr.meyn.search.plist ~/bin/fr.meyn.search.plist
	echo $(SEARCHPATH)
	sed -i '.orig' 's#SEARCHPATH#$(SEARCHPATH)#' ~/bin/fr.meyn.search.plist


smesearch: *.go cmd/*.go indexer/*.go searcher/*.go service/*.go
	go build smesearch.go


unload-plist:
	launchctl unload  ~/Library/LaunchAgents/fr.meyn.search.plist

install-plist: install
	cp $(SEARCHPATH)/fr.meyn.search.plist ~/Library/LaunchAgents
	launchctl unload ~/Library/LaunchAgents/fr.meyn.search.plist
	launchctl load ~/Library/LaunchAgents/fr.meyn.search.plist
	launchctl list | grep fr.meyn
