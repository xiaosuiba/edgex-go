/*******************************************************************************
 * Copyright 2017 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package metadata

import (
	"errors"
	"github.com/edgexfoundry/edgex-go/internal/pkg"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"

	types "github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
	"github.com/edgexfoundry/edgex-go/internal/core/metadata/operators/command"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
)

func restGetAllCommands(w http.ResponseWriter, _ *http.Request) {
	results, err := dbClient.GetAllCommands()
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(results) > Configuration.Service.MaxResultCount {
		LoggingClient.Error("Max limit exceeded")
		http.Error(w, errors.New("Max limit exceeded").Error(), http.StatusRequestEntityTooLarge)
		return
	}
	pkg.Encode(&results, w, LoggingClient)
}

func restGetCommandById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var did string = vars[ID]
	res, err := dbClient.GetCommandById(did)
	if err != nil {
		if err == db.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		LoggingClient.Error(err.Error())
		return
	}
	pkg.Encode(res, w, LoggingClient)
}

func restGetCommandByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := url.QueryUnescape(vars[NAME])
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	results, err := dbClient.GetCommandByName(n)
	if err != nil {
		if err == db.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		LoggingClient.Error(err.Error())
		return
	}
	pkg.Encode(results, w, LoggingClient)
}

func restGetCommandsByDeviceId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	did, err := url.QueryUnescape(vars[ID])
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	op := command.NewDeviceIdExecutor(dbClient, did)
	commands, err := op.Execute()
	if err != nil {
		LoggingClient.Error(err.Error())
		switch err.(type) {
		case *types.ErrItemNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	pkg.Encode(&commands, w, LoggingClient)
}
