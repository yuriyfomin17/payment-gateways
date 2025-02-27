package httpserver

import (
	"net/http"
)

func (h HttpServer) UpdateGatewayPriority(w http.ResponseWriter, r *http.Request) {
	var request GatewayPriorityUpdate
	err := DecodeRequest(r, &request)
	if err != nil {
		BadRequest("unsupported-content-type", err, w)
		return
	}
	if ValidateGatewayPriority(request.Priority) != nil {
		BadRequest("invalid-gateway-priority", err, w)
		return
	}
	err = h.gatewayService.UpdateGatewayPriority(r.Context(), request.GtId, request.Priority)
	if err != nil {
		BadRequest("gateway-priority-not-updated", err, w)
		return
	}
	RespondOK(w, "success", "application/json", map[string]any{
		"gt_id":    request.GtId,
		"priority": request.Priority,
	})
}
