#!/usr/bin/env python

version = '1.0.0'

# Copyright John Thornton 2015

import gtk
import os
import subprocess

class Buglump:

	def __init__(self):
		self.builder = gtk.Builder()
		self.builder.add_from_file('dxf.glade')
		self.builder.connect_signals(self)
		self.window = self.builder.get_object('main_window')
		self.aboutdialog = self.builder.get_object("aboutdialog")
		self.aboutdialog.set_version(version)
		self.tolerance = self.builder.get_object('tolerance_entry')
		self.status = self.builder.get_object("status_label")
		self.status.set_text('No File Open')
		self.current_folder = os.path.expanduser('~')
		self.label2 = self.builder.get_object('label2')
		self.file_name = ''
		self.window.show()
		self.ini_check()

	def on_window_destroy(self, object, data=None):
		gtk.main_quit()

	def on_file_quit(self, menuitem, data=None):
		gtk.main_quit()

	def on_file_open(self, menuitem, data=None):
		self.fcd = gtk.FileChooserDialog("Open...", None,
			gtk.FILE_CHOOSER_ACTION_OPEN,
			(gtk.STOCK_CANCEL, gtk.RESPONSE_CANCEL, gtk.STOCK_OPEN, gtk.RESPONSE_OK))
		if len(self.current_folder) > 0:
			self.fcd.set_current_folder(self.current_folder)
		self.response = self.fcd.run()
		if self.response == gtk.RESPONSE_OK:
			self.status.set_text('File Selected %s' % self.fcd.get_filename())
			self.file_name = "-f=" + self.fcd.get_filename()
			self.current_folder = os.path.dirname(self.fcd.get_uri()[7:])
		else:
			self.status.set_text('No File Open')
		self.fcd.destroy()

	def on_file_convert(self, file_name, data=None):
		if len(self.file_name) > 0:
			self.args = self.file_name
			self.result = subprocess.call("dxf2gcode %s" %self.args, shell=True)
			if self.result == 0:
				self.status.set_text('Processing Complete')
			else:
				self.status.set_text('Error %d Processing %s' % (self.result, self.file_name))
		else:
			self.status.set_text('No File Open')

	def on_view_test(self, item, data=None):
		pass

	def ini_check(self, data=None):
		ini_path = os.path.expanduser('~') + '/.config/dxf2emc'
		ini_file = ini_path + '/dxf2emc.ini'
		if not os.path.exists(ini_path):
			os.makedirs(ini_path, 0755)
		if not os.path.exists(ini_file):
			fo = open(ini_file,'w')
			fo.write('TOLERANCE=0.000001')
			fo.close
			message = 'Preferences File Created\nthis can be edited from Edit > Preferences'
			md = gtk.MessageDialog(self.window,
			gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO,
			gtk.BUTTONS_OK, message)
			md.run()
			md.destroy()
		self.tolerance.set_text('0.0000001')

	def on_help_about(self, menuitem, data=None):
		self.response = self.aboutdialog.run()
		self.aboutdialog.hide()

	def on_revert_prefrences(self, data=None):
		self.status.set_text('Prefrences Reverted not operational.')

	def on_save_preferences(self, data=None):
		self.status.set_text('Prefrences Saved not operational.')

if __name__ == '__main__':
  main = Buglump()
  gtk.main()
